// Package training は training モジュールの公開 Facade。
// 外部からは Module 経由でのみアクセス可能。
package training

import (
	"net/http"

	"github.com/Watari995/musclead/internal/shared/dbtx"
	traininghandler "github.com/Watari995/musclead/internal/training/internal/handler"
	traininginfra "github.com/Watari995/musclead/internal/training/internal/infra"
	trainingusecase "github.com/Watari995/musclead/internal/training/internal/usecase"
	"github.com/go-gorp/gorp/v3"
)

type Module struct {
	Handler http.Handler
}

func NewModule(dbmap *gorp.DbMap) *Module {
	// repo
	dbmap.AddTableWithName(traininginfra.TrainingModel{}, "trainings").SetKeys(false, "ID")
	dbmap.AddTableWithName(traininginfra.TrainingExerciseModel{}, "training_exercises").SetKeys(false, "ID")
	dbmap.AddTableWithName(traininginfra.TrainingSetModel{}, "training_sets").SetKeys(false, "ID")

	dbmap.AddTableWithName(traininginfra.ExerciseModel{}, "exercises").SetKeys(false, "ID")

	repo := traininginfra.NewTrainingRepository(dbmap)
	txManager := dbtx.NewTransactionManager(dbmap)

	// use-case
	findTraining := trainingusecase.NewFindTrainingByID(repo)
	listTrainings := trainingusecase.NewListTraining(repo)
	recordTraining := trainingusecase.NewRecordTraining(repo, txManager)
	updateTraining := trainingusecase.NewUpdateTraining(repo, txManager)
	deleteTraining := trainingusecase.NewDeleteTrainingByID(repo)

	return &Module{Handler: traininghandler.New(findTraining, listTrainings, recordTraining, updateTraining, deleteTraining)}
}
