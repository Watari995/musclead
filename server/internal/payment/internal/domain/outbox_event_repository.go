package paymentdomain

import (
	"context"
)

// OutboxEventRepository は outbox_events テーブルへの永続化を抽象化する。
//
// 設計 (ADR 0015):
//   - Webhook 処理の TX 内で Save (INSERT) し、 TX 外で SNS publish する
//   - 即時 publish 成功時は MarkPublished() → Save (UPDATE) で published_at を SET
//   - 即時 publish 失敗時は published_at = NULL のまま残り、 outbox-relay Lambda が拾う
type OutboxEventRepository interface {
	// FindPending は未配信 (published_at IS NULL) の outbox event を limit 件返す。
	// outbox-relay Lambda が定期実行で呼び、 SNS publish の failsafe に使う。
	// 古い順 (FIFO) で取得する想定。
	FindPending(ctx context.Context, limit int) ([]*OutboxEvent, error)

	// Save は event を保存する (INSERT / UPDATE 兼用)。
	// 新規 INSERT、 published_at の UPDATE、 publish_error の UPDATE すべてで使う。
	Save(ctx context.Context, event *OutboxEvent) error
}
