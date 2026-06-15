package publicfunctions

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

// SubscriptionQuery は purchase が他 module に公開する読み取り系 API。
//
// 設計 (Pro gate):
//   - feature module (training 等) が「この user は Pro か」 を問い合わせるために使う。
//   - 依存方向: feature module → purchase (一方向、 purchase は feature を知らない)。
type SubscriptionQuery interface {
	IsPro(ctx context.Context, userID valueobject.UserID) (bool, error)
}
