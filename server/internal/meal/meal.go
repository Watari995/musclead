// Package meal is the public facade of the meal module.
// Modular Monolith (strict) のため、 外部からは Module 経由でのみアクセス可能。
package meal

import (
	"net/http"

	mealhandler "github.com/Watari995/musclead/internal/meal/internal/handler"
	mealinfra "github.com/Watari995/musclead/internal/meal/internal/infra"
	mealusecase "github.com/Watari995/musclead/internal/meal/internal/usecase"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/go-gorp/gorp/v3"
)

// Module は meal モジュールの公開 API。
// HTTP ハンドラだけを外に出すことで、 内部の usecase / repository を隠蔽する。
type Module struct {
	Handler http.Handler
}

func NewModule(dbmap *gorp.DbMap, cdnBaseURL string) *Module {
	// repositoryを作成する
	dbmap.AddTableWithName(mealinfra.MealModel{}, "meals").SetKeys(false, "ID")
	dbmap.AddTableWithName(mealinfra.MealPhotoModel{}, "meal_photos").SetKeys(false, "ID")
	repo := mealinfra.NewMealRepository(dbmap)
	txManager := dbtx.NewTransactionManager(dbmap)

	record := mealusecase.NewRecordMeal(repo, txManager)
	find := mealusecase.NewFindMealByID(repo)
	update := mealusecase.NewUpdateMeal(repo, txManager)
	delete := mealusecase.NewDeleteMealByID(repo)
	list := mealusecase.NewListMeals(repo)

	return &Module{
		Handler: mealhandler.New(record, find, update, delete, list, cdnBaseURL),
	}
}
