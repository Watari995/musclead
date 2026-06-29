package userdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type UserPreferences struct {
	id            valueobject.UserPreferencesID
	userID        valueobject.UserID
	theme         valueobject.Theme
	mealColor     valueobject.ColorHex
	trainingColor valueobject.ColorHex
	weightColor   valueobject.ColorHex
	createdAt     time.Time
	updatedAt     time.Time
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

func (p *UserPreferences) MealColor() valueobject.ColorHex {
	return p.mealColor
}

func (p *UserPreferences) SetMealColor(color valueobject.ColorHex) {
	p.mealColor = color
	p.updatedAt = time.Now()
}

func (p *UserPreferences) TrainingColor() valueobject.ColorHex {
	return p.trainingColor
}

func (p *UserPreferences) SetTrainingColor(color valueobject.ColorHex) {
	p.trainingColor = color
	p.updatedAt = time.Now()
}

func (p *UserPreferences) WeightColor() valueobject.ColorHex {
	return p.weightColor
}

func (p *UserPreferences) SetWeightColor(color valueobject.ColorHex) {
	p.weightColor = color
	p.updatedAt = time.Now()
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
	trainingColor, err := valueobject.NewColorHex("#4A90E2")
	if err != nil {
		return nil
	}
	mealColor, err := valueobject.NewColorHex("#7ED321")
	if err != nil {
		return nil
	}
	weightColor, err := valueobject.NewColorHex("#FF6B6B") // Default weight color
	if err != nil {
		return nil
	}

	return NewUserPreferences(
		valueobject.NewPrimaryID[valueobject.UserPreferencesID](),
		userID,
		*theme,
		*trainingColor,
		*mealColor,
		*weightColor,
		time.Now(),
		time.Now(),
	)
}

func NewUserPreferences(
	id valueobject.UserPreferencesID,
	userID valueobject.UserID,
	theme valueobject.Theme,
	mealColor valueobject.ColorHex,
	trainingColor valueobject.ColorHex,
	weightColor valueobject.ColorHex,
	createdAt time.Time,
	updatedAt time.Time,
) *UserPreferences {
	return &UserPreferences{
		id:            id,
		userID:        userID,
		theme:         theme,
		mealColor:     mealColor,
		trainingColor: trainingColor,
		weightColor:   weightColor,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}
