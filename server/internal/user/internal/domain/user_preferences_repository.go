package userdomain

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

type UserPreferencesRepository interface {
	FindByUserID(ctx context.Context, userID valueobject.UserID) (*UserPreferences, error)
	Save(ctx context.Context, preferences *UserPreferences) error
}
