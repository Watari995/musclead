// Package training は training モジュールの公開 Facade。
// 外部からは Module 経由でのみアクセス可能。
package training

import (
	"net/http"
	"time"

	purchasepublicfunctions "github.com/Watari995/musclead/internal/purchase/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	trainingpublicfunctions "github.com/Watari995/musclead/internal/training/interface/publicfunctions"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	traininghandler "github.com/Watari995/musclead/internal/training/internal/handler"
	traininginfra "github.com/Watari995/musclead/internal/training/internal/infra"
	trainingusecase "github.com/Watari995/musclead/internal/training/internal/usecase"
	"github.com/go-gorp/gorp/v3"
	"github.com/redis/go-redis/v9"
)

type Module struct {
	TrainingHandler http.Handler
	ExerciseHandler http.Handler
	RoutineHandler  http.Handler
	trainingQuery   trainingpublicfunctions.TrainingQuery
}

// NewModule は training モジュールを初期化する。
// redisClient が nil の場合は NoOp キャッシュを使う（ローカル開発用）。
func NewModule(dbmap *gorp.DbMap, subscriptionQuery purchasepublicfunctions.SubscriptionQuery, redisClient *redis.Client) *Module {
	// == repo ==
	dbmap.AddTableWithName(traininginfra.TrainingModel{}, "trainings").SetKeys(false, "ID")
	dbmap.AddTableWithName(traininginfra.TrainingExerciseModel{}, "training_exercises").SetKeys(false, "ID")
	dbmap.AddTableWithName(traininginfra.TrainingSetModel{}, "training_sets").SetKeys(false, "ID")

	dbmap.AddTableWithName(traininginfra.ExerciseModel{}, "exercises").SetKeys(false, "ID")

	dbmap.AddTableWithName(traininginfra.RoutineExerciseModel{}, "routine_exercises").SetKeys(false, "ID")
	dbmap.AddTableWithName(traininginfra.RoutineModel{}, "routines").SetKeys(false, "ID")

	trainingRepo := traininginfra.NewTrainingRepository(dbmap)
	exerciseRepo := traininginfra.NewExerciseRepository(dbmap)
	routineRepo := traininginfra.NewRoutineRepository(dbmap)
	routineQueryService := traininginfra.NewRoutineQueryService(dbmap)
	exerciseRecordQueryService := traininginfra.NewExerciseRecordQueryService(dbmap)
	trainingQueryService := traininginfra.NewTrainingQueryService(dbmap)

	txManager := dbtx.NewTransactionManager(dbmap)

	// == cache ==
	// weight.NewModule と同じパターン。redisClient が nil のときは NoOp キャッシュを使う。
	var bestSetCache trainingdomain.ExerciseBestSetTimeseriesCache
	if redisClient != nil {
		bestSetCache = traininginfra.NewRedisExerciseBestSetTimeseriesCache(redisClient, 7*24*time.Hour)
	} else {
		bestSetCache = traininginfra.NewNoOpExerciseBestSetTimeseriesCache()
	}

	// == use-case ==
	// training
	findTraining := trainingusecase.NewFindTrainingByID(trainingRepo)
	listTrainings := trainingusecase.NewListTraining(trainingRepo)
	recordTraining := trainingusecase.NewRecordTraining(trainingRepo, txManager, bestSetCache)
	updateTraining := trainingusecase.NewUpdateTraining(trainingRepo, txManager, bestSetCache)
	deleteTraining := trainingusecase.NewDeleteTrainingByID(trainingRepo, bestSetCache)
	// exercise
	findExercise := trainingusecase.NewFindExerciseByID(exerciseRepo)
	listExercises := trainingusecase.NewListExercises(exerciseRepo)
	createExercise := trainingusecase.NewCreateExercise(exerciseRepo)
	updateExercise := trainingusecase.NewUpdateExercise(exerciseRepo)
	deleteExercise := trainingusecase.NewDeleteExerciseByID(exerciseRepo)
	reorderExercises := trainingusecase.NewReorderExercises(exerciseRepo, txManager)
	// routine
	findRoutine := trainingusecase.NewFindRoutineByID(routineQueryService)
	listRoutines := trainingusecase.NewListRoutines(routineQueryService)
	createRoutine := trainingusecase.NewCreateRoutine(routineRepo, subscriptionQuery)
	updateRoutine := trainingusecase.NewUpdateRoutine(routineRepo)
	deleteRoutine := trainingusecase.NewDeleteRoutineByID(routineRepo)
	reorderRoutines := trainingusecase.NewReorderRoutines(routineRepo, txManager)
	// exercise record
	findBestSets := trainingusecase.NewFindBestSetsByExerciseIDs(exerciseRecordQueryService)
	getBestSetTimeseries := trainingusecase.NewGetExerciseBestSetTimeseries(exerciseRecordQueryService, bestSetCache)
	findLastSessionSets := trainingusecase.NewFindLastSessionSetsByExerciseIDs(exerciseRecordQueryService)
	// calendar
	listTrainingDatesByMonth := trainingusecase.NewListTrainingDatesByMonth(trainingQueryService)
	listTrainingSummaryByDate := trainingusecase.NewListTrainingSummaryByDate(trainingQueryService)
	trainingQuery := trainingusecase.NewTrainingQuery(listTrainingDatesByMonth, listTrainingSummaryByDate)

	return &Module{
		TrainingHandler: traininghandler.NewTrainingHandler(findTraining, listTrainings, recordTraining, updateTraining, deleteTraining),
		ExerciseHandler: traininghandler.NewExerciseHandler(findExercise, findBestSets, getBestSetTimeseries, listExercises, findLastSessionSets, createExercise, updateExercise, deleteExercise, reorderExercises),
		RoutineHandler:  traininghandler.NewRoutineHandler(findRoutine, listRoutines, createRoutine, updateRoutine, deleteRoutine, reorderRoutines),
		trainingQuery:   trainingQuery,
	}
}

func (m *Module) TrainingQuery() trainingpublicfunctions.TrainingQuery {
	return m.trainingQuery
}
