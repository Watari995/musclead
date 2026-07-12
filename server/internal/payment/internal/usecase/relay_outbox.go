package paymentusecase

import (
	"context"
	"log/slog"

	"github.com/Watari995/musclead/internal/myerror"
	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
	shareddomain "github.com/Watari995/musclead/internal/shared/domain"
	userpublicfunctions "github.com/Watari995/musclead/internal/user/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/valueobject"
)

// relayBatchSize は1回のポーリングで処理する outbox の最大件数。
const relayBatchSize = 50

// RelayOutbox は outbox_events を拾って SQS に流す relay の「1回分」 の処理。
// worker goroutine が一定間隔でこの Execute を呼ぶ (ADR 0020 ①)。
//
// 流れ (1件ごと):
//  1. outboxRepo.FindPendingByEventTypes(payment 系 event type, relayBatchSize) で未配信を取得
//  2. payment_succeeded のみ: AggregateID(payment_id) → paymentRepo.FindByID → user_id
//     → userQuery.GetEmailByUserID で email を補完 → publisher.Publish
//  3. 種別を問わず MarkPublished → outboxRepo.Save (台帳から流す。 対象外の種別も溜めない)
//
// 冪等: ここは at-least-once (publish 後 Save 前に crash で再送あり)。 重複は consumer 側で吸収 (ADR 0020 ④⑤)。
type RelayOutbox struct {
	outboxRepo  shareddomain.OutboxEventRepository
	paymentRepo paymentdomain.PaymentRepository
	userQuery   userpublicfunctions.UserQuery
	publisher   paymentdomain.Publisher
}

func (uc *RelayOutbox) Execute(ctx context.Context) error {
	events, err := uc.outboxRepo.FindPendingByEventTypes(ctx, []valueobject.OutboxEventType{
		valueobject.NewOutboxEventTypeFromCode(valueobject.OutboxEventTypePaymentSucceeded),
		valueobject.NewOutboxEventTypeFromCode(valueobject.OutboxEventTypePaymentFailed),
		valueobject.NewOutboxEventTypeFromCode(valueobject.OutboxEventTypePaymentCanceled),
		valueobject.NewOutboxEventTypeFromCode(valueobject.OutboxEventTypePaymentRenewed),
	}, relayBatchSize)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	if len(events) == 0 {
		return nil // 処理するイベントがなかったら何もしない
	}
	for _, event := range events {
		// event_type = payment_succeededのみAggregateID(paymentId)で処理をする
		if event.EventType().IsPaymentSucceeded() {
			paymentID, err := valueobject.NewPrimaryIDFromString[valueobject.PaymentID](event.AggregateID())
			if err != nil {
				slog.Error("relay: failed to publish", "err", err)
				continue
			}
			payment, err := uc.paymentRepo.FindByID(ctx, *paymentID)
			if err != nil {
				slog.Error("relay: failed to publish", "err", err)
				continue
			}
			if payment == nil {
				continue
			}
			output, err := uc.userQuery.GetEmailByUserID(ctx, userpublicfunctions.GetEmailByUserIDInput{UserID: payment.UserID()})
			if err != nil {
				slog.Error("relay: failed to publish", "err", err)
				continue
			}
			if err := uc.publisher.Publish(ctx, paymentdomain.PublishMessage{
				EventID: event.ID().Value(),
				Type:    event.EventType().Value(),
				Email:   output.Email.Value(),
			}); err != nil {
				continue // あとで再度やるので何もしないでスルーする
			}
		}
		// 今後他のoutbox eventの処理が入る
		// publishして保存
		event.MarkPublished()
		if err := uc.outboxRepo.Save(ctx, event); err != nil {
			return myerror.NewInternalError().Wrap(err)
		}
	}
	return nil
}

func NewRelayOutbox(
	outboxRepo shareddomain.OutboxEventRepository,
	paymentRepo paymentdomain.PaymentRepository,
	userQuery userpublicfunctions.UserQuery,
	publisher paymentdomain.Publisher,
) *RelayOutbox {
	return &RelayOutbox{
		outboxRepo:  outboxRepo,
		paymentRepo: paymentRepo,
		userQuery:   userQuery,
		publisher:   publisher,
	}
}
