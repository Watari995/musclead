package healthsyncdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

type TokenRepository interface {
	FindByUserID(ctx context.Context, userID valueobject.UserID) (*Token, error)
	FindAllActive(ctx context.Context) ([]*Token, error)
	Save(ctx context.Context, token *Token) error
}
