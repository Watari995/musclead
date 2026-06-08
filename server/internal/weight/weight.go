// Package weight is the public facade of the weight module.
// Modular Monolith (strict) のため、 外部からは Module 経由でのみアクセス可能。
package weight

import (
	"net/http"

	weighthandler "github.com/Watari995/musclead/internal/weight/internal/handler"
	weightinfra "github.com/Watari995/musclead/internal/weight/internal/infra"
	weightusecase "github.com/Watari995/musclead/internal/weight/internal/usecase"
	"github.com/go-gorp/gorp/v3"
)

type Module struct {
	Handler http.Handler
}

func NewModule(dbmap *gorp.DbMap) *Module {
	// repository
	dbmap.AddTableWithName(weightinfra.WeightModel{}, "weights").SetKeys(false, "ID")
	repo := weightinfra.NewWeightRepository(dbmap)

	// use-case
	record := weightusecase.NewRecordWeight(repo)
	find := weightusecase.NewFindWeightByID(repo)
	list := weightusecase.NewListWeights(repo)
	update := weightusecase.NewUpdateWeight(repo)
	delete := weightusecase.NewDeleteWeightByID(repo)
	return &Module{
		Handler: weighthandler.New(record, find, list, update, delete),
	}
}
