// Package main is the entry point of the musclead API server.
//
//	@title			musclead API
//	@version		0.1.0
//	@description	musclead (筋トレ・食事・体重 一元管理 SaaS) のバックエンド API。
//	@host			localhost:8080
//	@BasePath		/
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				"Bearer <access_token>" 形式で指定する。
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Watari995/musclead/internal/auth"
	"github.com/Watari995/musclead/internal/billing"
	"github.com/Watari995/musclead/internal/calendar"
	"github.com/Watari995/musclead/internal/food"
	"github.com/Watari995/musclead/internal/healthsync"
	"github.com/Watari995/musclead/internal/meal"
	"github.com/Watari995/musclead/internal/notification"
	"github.com/Watari995/musclead/internal/payment"
	"github.com/Watari995/musclead/internal/purchase"
	_ "github.com/Watari995/musclead/internal/shared"
	shareddomain "github.com/Watari995/musclead/internal/shared/domain"
	"github.com/Watari995/musclead/internal/shared/httpx"
	cacheinfra "github.com/Watari995/musclead/internal/shared/infra/cache"
	sharedstorage "github.com/Watari995/musclead/internal/shared/infra/storage"
	"github.com/Watari995/musclead/internal/training"
	"github.com/Watari995/musclead/internal/user"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/Watari995/musclead/internal/weight"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-gorp/gorp/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	if err := run(); err != nil {
		slog.Error("server terminated with error", "error", err)
		os.Exit(1)
	}
}

// storagePublicBaseURL は画像公開 URL のベースを返す。
// STORAGE_PUBLIC_BASE_URL があればそれ(R2 の公開ドメイン等)、無ければ AWS S3 の形式を組み立てる。
func storagePublicBaseURL(bucket string) string {
	if base := os.Getenv("STORAGE_PUBLIC_BASE_URL"); base != "" {
		return base
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com", bucket, os.Getenv("AWS_REGION"))
}

func run() error {
	_ = godotenv.Load("../.env", ".env")

	addr := getenv("ADDR", ":8080")

	db, err := openDB()
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	dbmap := &gorp.DbMap{
		Db:      db,
		Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8MB4"},
	}
	// sentry の初期化
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		TracesSampleRate: 1.0,
	}); err != nil {
		return fmt.Errorf("sentry init: %w", err)
	}
	defer sentry.Flush(2 * time.Second)

	// S3 Clientを作成
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("aws config: %v", err)
	}
	// オブジェクトストレージは S3 互換 (AWS S3 / Cloudflare R2)。
	// STORAGE_ENDPOINT を指定すると R2 等のカスタムエンドポイントへ向く(未指定なら AWS S3)。
	bucket := os.Getenv("STORAGE_BUCKET")
	s3RawClient := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if endpoint := os.Getenv("STORAGE_ENDPOINT"); endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})
	storageClient := sharedstorage.NewS3Client(s3RawClient, bucket)
	urlBuilder := sharedstorage.NewS3URLBuilder(storagePublicBaseURL(bucket))
	sqsClient := sqs.NewFromConfig(awsCfg)
	redisClient := newRedisClient(context.Background())
	slog.Info("redis client initialized", "type", fmt.Sprintf("%T", redisClient))
	mux, paymentModule, healthSyncModule := newMux(dbmap, storageClient, urlBuilder, redisClient, sqsClient)

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// outbox relay worker を起動 (ctx Done で停止)。 SQS 未設定時は no-op。
	go paymentModule.RunRelay(ctx)

	// 体重同期のポーリングを起動
	go healthSyncModule.RunSync(ctx)

	errCh := make(chan error, 1)
	go func() {
		slog.Info("server starting", "addr", addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	case err := <-errCh:
		return fmt.Errorf("listen: %w", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	slog.Info("server stopped gracefully")
	return nil
}

func openDB() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true&loc=UTC",
		getenv("DB_USER", "musclead"),
		getenv("DB_PASSWORD", "musclead"),
		getenv("DB_HOST", "127.0.0.1"),
		getenv("DB_PORT", "3306"),
		getenv("DB_NAME", "musclead"),
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}

// newMux は全モジュールの HTTP ハンドラをマウントしたルーターを返す。
// 各モジュールは自身の Handler を Module.Handler として公開する。
func newMux(dbmap *gorp.DbMap, storageClient shareddomain.StorageClient, urlBuilder shareddomain.URLBuilder, redisClient *redis.Client, sqsClient *sqs.Client) (http.Handler, *payment.Module, *healthsync.Module) {
	mux := http.NewServeMux()

	// ヘルスチェック
	mux.HandleFunc("GET /health", healthHandler)

	// 各モジュールを組み立て、 そのハンドラをマウント
	// swaggerのマウント
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	userModule := user.NewModule(dbmap, storageClient, urlBuilder)
	authModule := auth.NewModule(dbmap, userModule.UserCommand())
	mealModule := meal.NewModule(dbmap, storageClient, urlBuilder)
	foodModule := food.NewModule(dbmap, &http.Client{Timeout: 10 * time.Second})
	weightModule := weight.NewModule(dbmap, redisClient)
	healthSyncModule := healthsync.NewModule(
		dbmap,
		&http.Client{Timeout: 10 * time.Second},
		os.Getenv("HEALTH_PLANET_CLIENT_ID"),
		os.Getenv("HEALTH_PLANET_CLIENT_SECRET"),
		os.Getenv("JWT_SECRET"),
		getenv("FRONTEND_URL", "http://localhost:3000"),
		weightModule.WeightCommand(),
		weightModule.WeightQuery(),
	)
	paymentModule := payment.NewModule(dbmap, payment.Config{
		StripeAPIKey:               os.Getenv("STRIPE_SECRET_KEY"),
		StripeSuccessURL:           os.Getenv("STRIPE_SUCCESS_URL"),
		StripeCancelURL:            os.Getenv("STRIPE_CANCEL_URL"),
		StripeWebhookSigningSecret: os.Getenv("STRIPE_WEBHOOK_SIGNING_SECRET"),
		StripePortalReturnURL:      os.Getenv("STRIPE_PORTAL_RETURN_URL"),
		SQSQueueURL:                os.Getenv("OUTBOX_QUEUE_URL"),
		ResendAPIKey:               os.Getenv("RESEND_API_KEY"),
		MailFromAddress:            os.Getenv("MAIL_FROM_ADDRESS"),
	}, userModule.UserQuery(), sqsClient)
	priceIDByPlan := map[valueobject.SubscriptionPlanCode]string{
		valueobject.SubscriptionPlanPro: os.Getenv("STRIPE_PRO_PRICE_ID"),
	}
	purchaseModule := purchase.NewModule(dbmap, paymentModule.Command(), userModule.UserQuery(), priceIDByPlan)
	billingModule := billing.NewModule(paymentModule.WebhookCommand(), paymentModule.Processor(), purchaseModule.PurchaseCommand())
	trainingModule := training.NewModule(dbmap, purchaseModule.SubscriptionQuery(), redisClient)
	calendarModule := calendar.NewModule(trainingModule.TrainingQuery(), mealModule.MealQuery(), weightModule.WeightQuery())
	notificationModule := notification.NewModule(dbmap)

	// users
	mux.Handle("/users", userModule.PublicHandler)
	mux.Handle("/users/", authModule.Middleware(userModule.Handler))
	// auth
	mux.Handle("/auth/", authModule.Handler)
	// meals
	mux.Handle("/meals", authModule.Middleware(mealModule.Handler))
	mux.Handle("/meals/", authModule.Middleware(mealModule.Handler))
	// meal_templates
	mux.Handle("/meal_templates", authModule.Middleware(mealModule.Handler))
	mux.Handle("/meal_templates/", authModule.Middleware(mealModule.Handler))
	// trainings
	mux.Handle("/trainings", authModule.Middleware(trainingModule.TrainingHandler))
	mux.Handle("/trainings/", authModule.Middleware(trainingModule.TrainingHandler))
	// exercises
	mux.Handle("/exercises", authModule.Middleware(trainingModule.ExerciseHandler))
	mux.Handle("/exercises/", authModule.Middleware(trainingModule.ExerciseHandler))
	// routines
	mux.Handle("/routines", authModule.Middleware(trainingModule.RoutineHandler))
	mux.Handle("/routines/", authModule.Middleware(trainingModule.RoutineHandler))
	// weights
	mux.Handle("/weights", authModule.Middleware(weightModule.Handler))
	mux.Handle("/weights/", authModule.Middleware(weightModule.Handler))
	// healthsync
	mux.Handle("/integrations/healthplanet/auth", authModule.Middleware(healthSyncModule.Handler))
	mux.Handle("/integrations/healthplanet/start", healthSyncModule.Handler)
	mux.Handle("/integrations/healthplanet/callback", healthSyncModule.Handler)
	// purchase
	mux.Handle("/purchase", purchaseModule.Handler)
	mux.Handle("/purchase/", authModule.Middleware(purchaseModule.Handler))
	// billing (Stripe Webhook、 auth middleware なし)
	mux.Handle("/billing/", billingModule.Handler)
	// food
	mux.Handle("/food_products", authModule.Middleware(foodModule.Handler))
	mux.Handle("/food_products/", authModule.Middleware(foodModule.Handler))
	// calendar
	mux.Handle("/calendar", authModule.Middleware(calendarModule.Handler))
	mux.Handle("/calendar/", authModule.Middleware(calendarModule.Handler))
	// notification
	mux.Handle("/notifications", authModule.Middleware(notificationModule.Handler))
	mux.Handle("/notifications/", authModule.Middleware(notificationModule.Handler))

	sentryHandler := sentryhttp.New(sentryhttp.Options{Repanic: true})
	return sentryHandler.Handle(httpx.CORSMiddleware(mux)), paymentModule, healthSyncModule
}

// healthHandler はサーバー稼働確認用のシンプルなヘルスチェック。
//
//	@Summary	ヘルスチェック
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/health [get]
func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func newRedisClient(ctx context.Context) *redis.Client {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		slog.Info("REDIS_HOST not set, using nil redis client")
		return nil
	}
	port := getenv("REDIS_PORT", "6379")
	addr := host + ":" + port
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  500 * time.Millisecond,
		ReadTimeout:  100 * time.Millisecond,
		WriteTimeout: 100 * time.Millisecond,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		slog.Warn("redis ping failed", "err", err, "addr", addr)
		_ = client.Close()
		return nil
	}
	return client
}

func newCache(client *redis.Client) shareddomain.Cache {
	if client == nil {
		return cacheinfra.NewNoOpCache()
	}
	return cacheinfra.NewRedisCache(client)
}
