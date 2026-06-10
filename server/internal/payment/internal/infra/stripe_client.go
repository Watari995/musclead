package paymentinfra

import (
	"context"
	"errors"

	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stripe/stripe-go/v82"
)

// errStripeClientNotImplemented は skeleton 状態でメソッドが呼ばれた時に返す sentinel error。
// User が各メソッドの中身を実装したら、 各 method 内で削除する。
var errStripeClientNotImplemented = errors.New("stripe client method not implemented")

// stripeClient は paymentdomain.StripeClient interface の実装。
//
// 設計:
//   - 環境差分 (API key, URL, signing secret) はコンストラクタで受け取る (ADR 0017)
//   - Stripe SDK は stripe.Key のグローバル設定で API key を渡す慣習に従う
//   - 商品差分 (PriceID 等) は各メソッドの input struct で受け取る
//
// 実装方針:
//   - usecase / domain は Stripe SDK の型 (stripe.Customer 等) を見ない (ACL)
//   - 必要な情報だけ domain DTO (CreateCheckoutSessionOutput / ParseWebhookEventOutput 等) に詰めて返す
type stripeClient struct {
	successURL           string // Checkout 成功時の戻り URL (環境差分、 env で設定)
	cancelURL            string // Checkout キャンセル時の戻り URL (環境差分)
	webhookSigningSecret string // Webhook 署名検証用 secret (環境差分、 SSM 保管)
	portalReturnURL      string // Customer Portal の戻り URL (環境差分)
}

// NewStripeClient は stripeClient を初期化する。
// Stripe SDK のグローバル変数 stripe.Key に API key を設定するため、
// プロセス内で 1 度だけ呼ぶ想定 (Composition Root から)。
func NewStripeClient(apiKey, successURL, cancelURL, webhookSigningSecret, portalReturnURL string) paymentdomain.StripeClient {
	stripe.Key = apiKey
	return &stripeClient{
		successURL:           successURL,
		cancelURL:            cancelURL,
		webhookSigningSecret: webhookSigningSecret,
		portalReturnURL:      portalReturnURL,
	}
}

// CreateCustomer は Stripe 側に新規 Customer を作成し、 customer_id (cus_xxx) を返す。
// userID は metadata の "musclead_user_id" に格納し、 Stripe Dashboard での監査・サポートで使う。
//
// TODO (User 実装):
//
//	import "github.com/stripe/stripe-go/v82/customer"
//
//	params := &stripe.CustomerParams{
//	    Email: stripe.String(input.Email.Value()),
//	    Metadata: map[string]string{
//	        "musclead_user_id": input.UserID.Value(),
//	    },
//	}
//	params.Context = ctx
//	cust, err := customer.New(params)
//	if err != nil { return "", err }
//	return cust.ID, nil
func (s *stripeClient) CreateCustomer(ctx context.Context, input paymentdomain.CreateCustomerInput) (string, error) {
	return "", errStripeClientNotImplemented
}

// CreateCheckoutSession は Stripe Checkout Session を作成し、 URL と SessionID を返す。
//
// 設計 (ADR 0014):
//   - PaymentID を Stripe Idempotency-Key として使う
//     → ネットワーク失敗で同じ呼び出しが 2 回起きても Session は 1 個しか作られない
//   - mode は "subscription" 固定 (サブスク決済)
//   - SuccessURL / CancelURL はコンストラクタで受け取った env から取る
//
// TODO (User 実装):
//
//	import checkoutsession "github.com/stripe/stripe-go/v82/checkout/session"
//
//	params := &stripe.CheckoutSessionParams{
//	    Customer:   stripe.String(input.CustomerID),
//	    Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
//	    SuccessURL: stripe.String(s.successURL),
//	    CancelURL:  stripe.String(s.cancelURL),
//	    LineItems: []*stripe.CheckoutSessionLineItemParams{
//	        {Price: stripe.String(input.PriceID), Quantity: stripe.Int64(1)},
//	    },
//	}
//	params.Context = ctx
//	params.SetIdempotencyKey(input.PaymentID.Value())
//	sess, err := checkoutsession.New(params)
//	if err != nil { return paymentdomain.CreateCheckoutSessionOutput{}, err }
//	checkoutURL, err := valueobject.NewURL(sess.URL)
//	if err != nil { return paymentdomain.CreateCheckoutSessionOutput{}, err }
//	return paymentdomain.CreateCheckoutSessionOutput{
//	    SessionID:          sess.ID,
//	    CheckoutSessionURL: *checkoutURL,
//	}, nil
func (s *stripeClient) CreateCheckoutSession(ctx context.Context, input paymentdomain.CreateCheckoutSessionInput) (paymentdomain.CreateCheckoutSessionOutput, error) {
	return paymentdomain.CreateCheckoutSessionOutput{}, errStripeClientNotImplemented
}

// ParseWebhookEvent は Stripe Webhook の署名検証 + event 取り出しを行う。
//
// 設計 (ADR 0018):
//   - 署名検証は HMAC-SHA256 (webhook signing secret 必須)
//   - 署名 NG / パース失敗時は error → handler が 401 で弾く (改ざん / なりすまし対策)
//   - event.Data.Raw を JSON unmarshal して valueobject.Metadata に詰める
//   - usecase は EventType でブランチし、 Payload (map) から必要な field を取り出す
//
// TODO (User 実装):
//
//	import (
//	    "encoding/json"
//	    "github.com/stripe/stripe-go/v82/webhook"
//	)
//
//	event, err := webhook.ConstructEvent(input.Payload, input.SignatureHeader, s.webhookSigningSecret)
//	if err != nil { return paymentdomain.ParseWebhookEventOutput{}, err }
//
//	var payload valueobject.Metadata
//	if err := json.Unmarshal(event.Data.Raw, &payload); err != nil {
//	    return paymentdomain.ParseWebhookEventOutput{}, err
//	}
//
//	return paymentdomain.ParseWebhookEventOutput{
//	    StripeEventID: event.ID,
//	    EventType:     string(event.Type),
//	    Payload:       payload,
//	}, nil
func (s *stripeClient) ParseWebhookEvent(ctx context.Context, input paymentdomain.ParseWebhookEventInput) (paymentdomain.ParseWebhookEventOutput, error) {
	return paymentdomain.ParseWebhookEventOutput{}, errStripeClientNotImplemented
}

// CreatePortalSession は Stripe Customer Portal Session を作成し、 リダイレクト URL を返す。
// 解約 / カード変更フローで、 クライアントは返り値の URL に遷移する (window.location = portalURL)。
//
// 設計 (ADR 0017):
//   - Customer Portal は Stripe Dashboard で「キャンセル可」「決済手段更新可」 を設定済みの前提
//   - ReturnURL はコンストラクタで受け取った env を使う (環境差分)
//
// TODO (User 実装):
//
//	import billingportalsession "github.com/stripe/stripe-go/v82/billingportal/session"
//
//	params := &stripe.BillingPortalSessionParams{
//	    Customer:  stripe.String(customerID),
//	    ReturnURL: stripe.String(s.portalReturnURL),
//	}
//	params.Context = ctx
//	sess, err := billingportalsession.New(params)
//	if err != nil { return valueobject.URL{}, err }
//	portalURL, err := valueobject.NewURL(sess.URL)
//	if err != nil { return valueobject.URL{}, err }
//	return *portalURL, nil
func (s *stripeClient) CreatePortalSession(ctx context.Context, customerID string) (valueobject.URL, error) {
	return valueobject.URL{}, errStripeClientNotImplemented
}
