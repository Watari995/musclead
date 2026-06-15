package paymentinfra

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/stripe/stripe-go/v82"
	billingportalsession "github.com/stripe/stripe-go/v82/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/subscription"
	"github.com/stripe/stripe-go/v82/webhook"
)

// stripeClient は paymentdomain.StripeClient interface の実装。
//
// 設計:
//   - 環境差分 (API key, URL, signing secret) はコンストラクタで受け取る (ADR 0017)
//   - Stripe SDK は stripe.Key のグローバル設定で API key を渡す慣習に従う
//   - 商品差分 (PriceID 等) は各メソッドの input struct で受け取る
//   - usecase / domain は Stripe SDK の型 (stripe.Customer 等) を見ない (ACL)、
//     必要な情報だけ domain DTO に詰めて返す
type stripeClient struct {
	successURL           string // Checkout 成功時の戻り URL
	cancelURL            string // Checkout キャンセル時の戻り URL
	webhookSigningSecret string // Webhook 署名検証用 secret
	portalReturnURL      string // Customer Portal の戻り URL
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
func (s *stripeClient) CreateCustomer(ctx context.Context, input paymentdomain.CreateCustomerInput) (string, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(input.Email.Value()),
		Metadata: map[string]string{
			"musclead_user_id": input.UserID.Value(),
		},
	}
	params.Context = ctx
	cust, err := customer.New(params)
	if err != nil {
		return "", err
	}
	return cust.ID, nil
}

// CreateCheckoutSession は Stripe Checkout Session を作成し、 URL と SessionID を返す。
//
// 設計 (ADR 0014):
//   - PaymentID を Stripe Idempotency-Key として使う
//     → ネットワーク失敗で同じ呼び出しが 2 回起きても Session は 1 個しか作られない
//   - mode は "subscription" 固定
func (s *stripeClient) CreateCheckoutSession(ctx context.Context, input paymentdomain.CreateCheckoutSessionInput) (paymentdomain.CreateCheckoutSessionOutput, error) {
	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(input.CustomerID),
		// ClientReferenceID に PaymentID を載せると checkout.session.completed の
		// payload["client_reference_id"] で返ってくる。 これを使って Webhook 受信時に
		// 「どの payment か」 を特定する (ADR 0014 / X-2)。
		// subscription_id は InitiatePayment 時点では未確定なので引き当てキーに使えない。
		ClientReferenceID: stripe.String(input.PaymentID.Value()),
		Mode:              stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL:        stripe.String(s.successURL),
		CancelURL:         stripe.String(s.cancelURL),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{Price: stripe.String(input.PriceID), Quantity: stripe.Int64(1)},
		},
	}
	params.Context = ctx
	params.SetIdempotencyKey(input.PaymentID.Value())
	sess, err := checkoutsession.New(params)
	if err != nil {
		return paymentdomain.CreateCheckoutSessionOutput{}, err
	}
	checkoutURL, err := valueobject.NewURL(sess.URL)
	if err != nil {
		return paymentdomain.CreateCheckoutSessionOutput{}, err
	}
	return paymentdomain.CreateCheckoutSessionOutput{
		SessionID:          sess.ID,
		CheckoutSessionURL: *checkoutURL,
	}, nil
}

// ParseWebhookEvent は Stripe Webhook の署名検証 + event 取り出しを行う。
//
// 設計 (ADR 0018):
//   - 署名検証は HMAC-SHA256、 失敗時は error を返す (改ざん / なりすまし対策)
//   - event.Data.Raw を JSON unmarshal して valueobject.Metadata に詰める
//   - usecase は EventType でブランチし、 Payload (map) から必要な field を取り出す
func (s *stripeClient) ParseWebhookEvent(ctx context.Context, input paymentdomain.ParseWebhookEventInput) (paymentdomain.ParseWebhookEventOutput, error) {
	// IgnoreAPIVersionMismatch: Stripe アカウントの API version (例: 2026-05-27.dahlia) と
	// stripe-go SDK が期待する version (basil) が異なると ConstructEvent は reject する。
	// 我々は event.Data.Raw を自前で json.Unmarshal して map に詰めており、 SDK のオブジェクト型
	// デシリアライズに依存しないため、 version 差分を無視して安全に署名検証のみ行う。
	event, err := webhook.ConstructEventWithOptions(input.Payload, input.SignatureHeader, s.webhookSigningSecret,
		webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true})
	if err != nil {
		return paymentdomain.ParseWebhookEventOutput{}, err
	}
	var payload valueobject.Metadata
	if err := json.Unmarshal(event.Data.Raw, &payload); err != nil {
		return paymentdomain.ParseWebhookEventOutput{}, err
	}
	return paymentdomain.ParseWebhookEventOutput{
		StripeEventID: event.ID,
		EventType:     string(event.Type),
		Payload:       payload,
	}, nil
}

// CreatePortalSession は Stripe Customer Portal Session を作成し、 リダイレクト URL を返す。
// 解約 / カード変更フローで、 クライアントは返り値の URL に遷移する (window.location = portalURL)。
// ReturnURL はコンストラクタで受け取った env を使う (環境差分)。
func (s *stripeClient) CreatePortalSession(ctx context.Context, customerID string) (valueobject.URL, error) {
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(customerID),
		ReturnURL: stripe.String(s.portalReturnURL),
	}
	params.Context = ctx
	sess, err := billingportalsession.New(params)
	if err != nil {
		return valueobject.URL{}, err
	}
	portalURL, err := valueobject.NewURL(sess.URL)
	if err != nil {
		return valueobject.URL{}, err
	}
	return *portalURL, nil
}

// RetrieveSubscription は Stripe からサブスクリプションを取得し、 現在の課金期間終了時刻を返す。
//
// 設計:
//   - current_period_end は subscription 直下ではなく items.data[].current_period_end にある
//     (stripe-go v82 / Stripe API dahlia で確認済み。 最近の API 変更で item 単位へ移動)
//   - 単一プランなので item は 1 つ。 念のため空チェックして Data[0] を読む
func (s *stripeClient) RetrieveSubscription(ctx context.Context, subscriptionID string) (time.Time, error) {
	params := &stripe.SubscriptionParams{}
	params.Context = ctx
	sub, err := subscription.Get(subscriptionID, params)
	if err != nil {
		return time.Time{}, err
	}
	if sub.Items == nil || len(sub.Items.Data) == 0 {
		return time.Time{}, fmt.Errorf("subscription %s has no items", subscriptionID)
	}
	return time.Unix(sub.Items.Data[0].CurrentPeriodEnd, 0), nil
}
