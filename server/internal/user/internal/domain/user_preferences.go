package userdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type UserPreferences struct {
	id        valueobject.UserPreferencesID
	userID    valueobject.UserID
	theme     valueobject.Theme
	createdAt time.Time
	updatedAt time.Time
}

func (p *UserPreferences) ID() valueobject.UserPreferencesID {
	return p.id
}

func (p *UserPreferences) UserID() valueobject.UserID {
	return p.userID
}

func (p *UserPreferences) Theme() valueobject.Theme {
	return p.theme
}

func (p *UserPreferences) SetTheme(t valueobject.Theme) {
	p.theme = t
	p.updatedAt = time.Now()
}

func (p *UserPreferences) CreatedAt() time.Time {
	return p.createdAt
}

func (p *UserPreferences) UpdatedAt() time.Time {
	return p.updatedAt
}

func CreateDefaultUserPreferences(userID valueobject.UserID) *UserPreferences {
	theme, err := valueobject.NewTheme(string(valueobject.ThemeSystem))
	if err != nil {
		return nil
	}
	return NewUserPreferences(
		valueobject.NewPrimaryID[valueobject.UserPreferencesID](),
		userID,
		*theme,
		time.Now(),
		time.Now(),
	)
}

func NewUserPreferences(
	id valueobject.UserPreferencesID,
	userID valueobject.UserID,
	theme valueobject.Theme,
	createdAt time.Time,
	updatedAt time.Time,
) *UserPreferences {
	return &UserPreferences{
		id:        id,
		userID:    userID,
		theme:     theme,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
