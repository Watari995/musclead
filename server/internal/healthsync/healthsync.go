package healthsync

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	healthsynchandler "github.com/Watari995/musclead/internal/healthsync/internal/handler"
	healthsyncusecase "github.com/Watari995/musclead/internal/healthsync/internal/usecase"
	"github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
	"github.com/go-gorp/gorp/v3"

	healthsyncinfra "github.com/Watari995/musclead/internal/healthsync/internal/infra"
)

const syncInterval = 90 * time.Second

type Module struct {
	syncWeights *healthsyncusecase.SyncWeights
	Handler     http.Handler
}

func NewModule(
	dbmap *gorp.DbMap,
	httpClint *http.Client,
	clientID string,
	clientSecret string,
	jwtSecret string,
	frontendURL string,
	weightCommand publicfunctions.WeightCommand,
	weightQuery publicfunctions.WeightQuery,
) *Module {
	// repo
	dbmap.AddTableWithName(healthsyncinfra.HealthPlanetTokenModel{}, "healthplanet_tokens").SetKeys(false, "id")
	tokenRepo := healthsyncinfra.NewTokenRepository(dbmap)
	client := healthsyncinfra.NewHealthPlanetClient(httpClint, clientID, clientSecret)
	stateSigner := healthsyncinfra.NewJWTStateSigner(jwtSecret)

	// use-case
	connect := healthsyncusecase.NewConnectHealthPlanet(tokenRepo, client, stateSigner)
	buildAuthURL := healthsyncusecase.NewBuildAuthURL(stateSigner)
	sync := healthsyncusecase.NewSyncWeights(tokenRepo, client, client, weightCommand, weightQuery)

	healthsyncHandler := healthsynchandler.New(buildAuthURL, connect, clientID, frontendURL)

	return &Module{
		syncWeights: sync,
		Handler:     healthsyncHandler,
	}
}

// 90秒ごとに体重を同期する。
func (m *Module) RunSync(ctx context.Context) {
	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := m.syncWeights.Execute(ctx); err != nil {
				slog.Error("healthsync: failed to sync weights", "err", err)
			}
		}
	}
}
