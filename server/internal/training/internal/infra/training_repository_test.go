package traininginfra_test

import (
	"context"
	"testing"
	"time"

	"github.com/Watari995/musclead/internal/testhelper"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	traininginfra "github.com/Watari995/musclead/internal/training/internal/infra"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/shopspring/decimal"
)

func TestTrainingRepository_SaveAndFind(t *testing.T) {
	dbmap := testhelper.NewTestDB(t)

	dbmap.AddTableWithName(traininginfra.TrainingModel{}, "trainings").SetKeys(false, "ID")
	dbmap.AddTableWithName(traininginfra.ExerciseModel{}, "exercises").SetKeys(false, "ID")
	dbmap.AddTableWithName(traininginfra.TrainingSetModel{}, "training_sets").SetKeys(false, "ID")
	dbmap.AddTableWithName(traininginfra.TrainingExerciseModel{}, "training_exercises").SetKeys(false, "ID")

	ctx := context.Background()

	// 前提データ: ユーザー
	userID := valueobject.NewPrimaryID[valueobject.UserID]()
	userIDBytes, _ := userID.Bytes()
	_, err := dbmap.Db.ExecContext(ctx,
		`INSERT INTO users (id, name, email, password_hash, created_at, updated_at) VALUES (?, 'testuser', 'test@example.com', 'hash', NOW(6), NOW(6))`,
		userIDBytes,
	)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}

	// 前提データ: エクササイズ
	exerciseID := valueobject.NewPrimaryID[valueobject.ExerciseID]()
	exerciseIDBytes, _ := exerciseID.Bytes()
	_, err = dbmap.Db.ExecContext(ctx,
		`INSERT INTO exercises (id, user_id, name, display_order, created_at, updated_at) VALUES (?, ?,'Bench Press', 1, NOW(6), NOW(6))`,
		exerciseIDBytes, userIDBytes,
	)
	if err != nil {
		t.Fatalf("insert exercise: %v", err)
	}

	// Trainingドメインオブジェクトを作る
	weightKg, _ := valueobject.NewNonNegativeDecimal(decimal.NewFromFloat(100.0))
	setNumber, _ := valueobject.NewNonNegativeInt(1)
	reps, _ := valueobject.NewNonNegativeInt(10)
	displayOrder, _ := valueobject.NewNonNegativeInt(1)

	startedAt := time.Now().UTC().Truncate(time.Microsecond)

	training := trainingdomain.CreateTraining(trainingdomain.TrainingSpec{
		StartedAt: startedAt,
		EndedAt:   nil,
		Memo:      nil,
		Exercises: []trainingdomain.TrainingExerciseSpec{
			{
				ExerciseID:   exerciseID,
				DisplayOrder: *displayOrder,
				RestSeconds:  nil,
				Memo:         nil,
				Sets: []trainingdomain.TrainingSetSpec{
					{
						SetNumber:   *setNumber,
						WeightKg:    *weightKg,
						Reps:        *reps,
						RestSeconds: nil,
						Memo:        nil,
					},
				},
			},
		},
	}, userID)

	// Repositoryを作ってsave
	repo := traininginfra.NewTrainingRepository(dbmap)
	if err := repo.Save(ctx, training); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// find

	found, err := repo.FindByIDAndUserID(ctx, training.ID(), userID)
	if err != nil {
		t.Fatalf("FindByIDAndUserID: %v", err)
	}

	if found.ID() != training.ID() {
		t.Errorf("ID: want %v, got %v", training.ID(), found.ID())
	}

	if !found.StartedAt().Equal(training.StartedAt()) {
		t.Errorf("StartedAt: want %v, got %v", training.StartedAt(), found.StartedAt())
	}
	if len(found.Exercises()) != len(training.Exercises()) {
		t.Errorf("Exercises: want %v, got %v", len(training.Exercises()), len(found.Exercises()))
	}

	if len(found.Exercises()[0].Sets()) != len(training.Exercises()[0].Sets()) {
		t.Errorf("Sets: want %v, got %v", len(training.Exercises()[0].Sets()), len(found.Exercises()[0].Sets()))
	}

	gotWeight := found.Exercises()[0].Sets()[0].WeightKg().Value()
	wantWeight := training.Exercises()[0].Sets()[0].WeightKg().Value()

	if !gotWeight.Equal(wantWeight) {
		t.Errorf("WeightKg: want %v, got %v", wantWeight, gotWeight)
	}
}
