package mealusecase_test

import (
	"testing"
	"time"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

// テスト共通のダミー Meal 生成(指定 userId で作る)。
func newDummyMeal(t *testing.T, userID valueobject.UserID) *mealdomain.Meal {
	t.Helper()
	mealType, _ := valueobject.NewString20("lunch")
	calories, _ := valueobject.NewNonNegativeInt(600)
	return mealdomain.CreateMeal(
		userID,
		time.Now(),
		*mealType,
		*calories,
		nil, nil, nil, nil, nil,
	)
}
