// Package weight is the public facade of the weight module.
// Modular Monolith (strict) のため、 外部からは Module 経由でのみアクセス可能。
package weight

import (
	"net/http"
	"time"

	"github.com/Watari995/musclead/internal/weight/interface/publicfunctions"
	weightdomain "github.com/Watari995/musclead/internal/weight/internal/domain"
	weighthandler "github.com/Watari995/musclead/internal/weight/internal/handler"
	weightinfra "github.com/Watari995/musclead/internal/weight/internal/infra"
	weightusecase "github.com/Watari995/musclead/internal/weight/internal/usecase"
	"github.com/go-gorp/gorp/v3"
	"github.com/redis/go-redis/v9"
)

type Module struct {
	Handler       http.Handler
	weightCommand publicfunctions.WeightCommand
	weightQuery   publicfunctions.WeightQuery
}

func NewModule(dbmap *gorp.DbMap, redisClient *redis.Client) *Module {
	// repository
	dbmap.AddTableWithName(weightinfra.WeightModel{}, "weights").SetKeys(false, "ID")
	repo := weightinfra.NewWeightRepository(dbmap)
	var cache weightdomain.WeightTimeseriesCache
	if redisClient != nil {
		cache = weightinfra.NewRedisWeightTimeseriesCache(redisClient, 24*time.Hour)
	} else {
		cache = weightinfra.NewNoOpWeightTimeseriesCache()
	}

	// use-case
	record := weightusecase.NewRecordWeight(repo, cache)
	find := weightusecase.NewFindWeightByID(repo)
	list := weightusecase.NewListWeights(repo)
	update := weightusecase.NewUpdateWeight(repo, cache)
	delete := weightusecase.NewDeleteWeightByID(repo, cache)
	getTimeseries := weightusecase.NewGetWeightTimeseries(repo, cache)
	checkIfExistsWeightByUserIDAndMeasuredAt := weightusecase.NewCheckIfExistsWeightByUserIDAndMeasuredAt(repo)
	weightCommand := weightusecase.NewWeightCommand(record)
	weightQuery := weightusecase.NewWeightQuery(checkIfExistsWeightByUserIDAndMeasuredAt)

	return &Module{
		Handler:       weighthandler.New(record, find, list, update, delete, getTimeseries),
		weightCommand: weightCommand,
		weightQuery:   weightQuery,
	}
}

func (m *Module) WeightCommand() publicfunctions.WeightCommand {
	return m.weightCommand
}

func (m *Module) WeightQuery() publicfunctions.WeightQuery {
	return m.weightQuery
}
