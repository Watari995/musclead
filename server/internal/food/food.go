// Package food は食品マスタモジュールの公開インターフェース。
// バーコード検索・名前検索・ユーザー登録を提供する。
package food

import (
	"net/http"

	foodhandler "github.com/Watari995/musclead/internal/food/internal/handler"
	foodinfra "github.com/Watari995/musclead/internal/food/internal/infra"
	foodusecase "github.com/Watari995/musclead/internal/food/internal/usecase"
	"github.com/go-gorp/gorp/v3"
)

// Module は food モジュールの公開 API。
type Module struct {
	Handler http.Handler
}

func NewModule(dbmap *gorp.DbMap, openFoodFactsClient *http.Client) *Module {
	// == repo ==
	dbmap.AddTableWithName(foodinfra.FoodProductModel{}, "food_products").SetKeys(false, "ID")
	foodProductRepo := foodinfra.NewFoodProductRepository(dbmap)
	foodFactClient := foodinfra.NewOpenFoodFactsClient(openFoodFactsClient)

	// == use-case ==
	searchByName := foodusecase.NewSearchByName(foodProductRepo)
	searchByBarcode := foodusecase.NewSearchByBarcode(foodProductRepo, foodFactClient)
	createFoodProduct := foodusecase.NewCreateFoodProduct(foodProductRepo)

	return &Module{Handler: foodhandler.New(searchByName, searchByBarcode, createFoodProduct)}
}
