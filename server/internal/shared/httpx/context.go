package httpx

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
)

type contextKey string

// contextKeyはcontextのキー衝突を防ぐ
const userIDKey contextKey = "userId"

// middlewareが userIDを載せる時に使う
func WithUserID(ctx context.Context, userID valueobject.UserID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (valueobject.UserID, error) {
	userID, ok := ctx.Value(userIDKey).(valueobject.UserID)
	if !ok {
		return valueobject.UserID{}, myerror.NewUnauthorizedError().SetMessage("user id not found in context")
	}
	return userID, nil
}
