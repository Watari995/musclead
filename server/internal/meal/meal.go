// Package meal is the public facade of the meal module.
// Modular Monolith (strict) のため、 外部からは Module 経由でのみアクセス可能。
package meal

import (
	"net/http"

	mealhandler "github.com/Watari995/musclead/internal/meal/internal/handler"
	mealinfra "github.com/Watari995/musclead/internal/meal/internal/infra"
	mealusecase "github.com/Watari995/musclead/internal/meal/internal/usecase"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	shareddomain "github.com/Watari995/musclead/internal/shared/domain"
	"github.com/go-gorp/gorp/v3"
)

// Module は meal モジュールの公開 API。
// HTTP ハンドラだけを外に出すことで、 内部の usecase / repository を隠蔽する。
type Module struct {
	Handler http.Handler
}

func NewModule(dbmap *gorp.DbMap, storageClient shareddomain.StorageClient, urlBuilder shareddomain.URLBuilder) *Module {
	dbmap.AddTableWithName(mealinfra.MealModel{}, "meals").SetKeys(false, "ID")
	dbmap.AddTableWithName(mealinfra.MealPhotoModel{}, "meal_photos").SetKeys(false, "ID")
	dbmap.AddTableWithName(mealinfra.MealTemplateModel{}, "meal_templates").SetKeys(false, "ID")

	repo := mealinfra.NewMealRepository(dbmap)
	templateRepo := mealinfra.NewMealTemplateRepository(dbmap)
	txManager := dbtx.NewTransactionManager(dbmap)

	record := mealusecase.NewRecordMeal(repo, txManager)
	find := mealusecase.NewFindMealByID(repo)
	update := mealusecase.NewUpdateMeal(repo, txManager, storageClient)
	delete := mealusecase.NewDeleteMealByID(repo, storageClient)
	list := mealusecase.NewListMeals(repo)
	generateMealPhotoImagePresignedURL := mealusecase.NewGenerateMealPhotoImagePresignedURL(storageClient)

	createTemplate := mealusecase.NewCreateMealTemplate(templateRepo)
	listTemplates := mealusecase.NewListMealTemplates(templateRepo)
	updateTemplate := mealusecase.NewUpdateMealTemplate(templateRepo)
	deleteTemplate := mealusecase.NewDeleteMealTemplate(templateRepo)
	reorderTemplate := mealusecase.NewReorderMealTemplate(templateRepo, txManager)

	return &Module{
		Handler: mealhandler.New(
			urlBuilder,
			record, find, update, delete, list, generateMealPhotoImagePresignedURL,
			listTemplates, createTemplate, updateTemplate, deleteTemplate, reorderTemplate,
		),
	}
}
