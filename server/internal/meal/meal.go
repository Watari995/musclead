// Package meal is the public facade of the meal module.
// Modular Monolith (strict) のため、 外部からは Module 経由でのみアクセス可能。
package meal

import (
	"database/sql"
	"net/http"

	mealhandler "github.com/Watari995/musclead/internal/meal/internal/handler"
	mealinfra "github.com/Watari995/musclead/internal/meal/internal/infra"
	mealusecase "github.com/Watari995/musclead/internal/meal/internal/usecase"
)

// Module は meal モジュールの公開 API。
// HTTP ハンドラだけを外に出すことで、 内部の usecase / repository を隠蔽する。
type Module struct {
	Handler http.Handler
}

func NewModule(db *sql.DB) *Module {
	repo := mealinfra.NewMealRepository(db)

	listMeals := mealusecase.NewListMeals(repo)
	deleteMealByID := mealusecase.NewDeleteMealByID(repo)

	return &Module{
		Handler: mealhandler.New(listMeals, deleteMealByID),
	}
}
