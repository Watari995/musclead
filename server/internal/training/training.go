// Package training は training モジュールの公開 Facade。
// 外部からは Module 経由でのみアクセス可能。
package training

import (
	"net/http"

	purchasepublicfunctions "github.com/Watari995/musclead/internal/purchase/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	traininghandler "github.com/Watari995/musclead/internal/training/internal/handler"
	traininginfra "github.com/Watari995/musclead/internal/training/internal/infra"
	trainingusecase "github.com/Watari995/musclead/internal/training/internal/usecase"
	"github.com/go-gorp/gorp/v3"
)

type Module struct {
	TrainingHandler http.Handler
	ExerciseHandler http.Handler
	RoutineHandler  http.Handler
}

func NewModule(dbmap *gorp.DbMap, subscriptionQuery purchasepublicfunctions.SubscriptionQuery) *Module {
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

	txManager := dbtx.NewTransactionManager(dbmap)

	// == use-case ==
	// training
	findTraining := trainingusecase.NewFindTrainingByID(trainingRepo)
	listTrainings := trainingusecase.NewListTraining(trainingRepo)
	recordTraining := trainingusecase.NewRecordTraining(trainingRepo, txManager)
	updateTraining := trainingusecase.NewUpdateTraining(trainingRepo, txManager)
	deleteTraining := trainingusecase.NewDeleteTrainingByID(trainingRepo)
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
	// exercise record
	findBestSets := trainingusecase.NewFindBestSetsByExerciseIDs(exerciseRecordQueryService)

	return &Module{TrainingHandler: traininghandler.NewTrainingHandler(findTraining, listTrainings, recordTraining, updateTraining, deleteTraining), ExerciseHandler: traininghandler.NewExerciseHandler(findExercise, findBestSets, listExercises, createExercise, updateExercise, deleteExercise, reorderExercises), RoutineHandler: traininghandler.NewRoutineHandler(findRoutine, listRoutines, createRoutine, updateRoutine, deleteRoutine)}
}
