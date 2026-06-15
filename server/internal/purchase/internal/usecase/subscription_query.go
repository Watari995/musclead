package purchaseusecase

import (
	"context"

	purchasepublicfunctions "github.com/Watari995/musclead/internal/purchase/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/valueobject"
)

// subscriptionQuery は publicfunctions.SubscriptionQuery を満たす adapter。
// Pro gate (feature module) からの is_pro 問い合わせ用。
// is_pro 判定は GetSubscription を再利用する (判定ロジックを1箇所に保つ)。
type subscriptionQuery struct {
	getSubscription *GetSubscription
}

func NewSubscriptionQuery(getSubscription *GetSubscription) purchasepublicfunctions.SubscriptionQuery {
	return &subscriptionQuery{getSubscription: getSubscription}
}

func (q *subscriptionQuery) IsPro(ctx context.Context, userID valueobject.UserID) (bool, error) {
	output, err := q.getSubscription.Execute(ctx, GetSubscriptionInput{UserID: userID})
	if err != nil {
		return false, err
	}
	return output.IsPro, nil
}
