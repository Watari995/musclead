package shareddomain

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

// OutboxEventRepository は outbox_events テーブルへの永続化を抽象化する。
type OutboxEventRepository interface {
	// FindPending は未配信 (published_at IS NULL) の outbox event を limit 件返す。
	// 古い順 (FIFO) で取得する想定。
	FindPendingByEventTypes(ctx context.Context, eventTypes []valueobject.OutboxEventType, limit int) ([]*OutboxEvent, error)

	// Save は event を保存する (INSERT / UPDATE 兼用)。
	// 新規 INSERT、 published_at の UPDATE、 publish_error の UPDATE すべてで使う。
	Save(ctx context.Context, event *OutboxEvent) error
}
