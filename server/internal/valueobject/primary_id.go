package valueobject

type PrimaryID interface {
	// PrimaryIDBase を埋め込んだ型を指定
	~struct{ PrimaryIDBase }
}

func NewPrimaryID[T PrimaryID]() T {
	return T{PrimaryIDBase: newPrimaryIDBase()}
}

func NewPrimaryIDFromString[T PrimaryID](s string) (*T, error) {
	p, err := newPrimaryIDBaseFromString(s)
	if err != nil {
		return nil, err
	}
	return &T{PrimaryIDBase: p}, nil
}

// ID 型一覧(Shared Kernel として全モジュール共通の参照子)
type SessionID struct{ PrimaryIDBase }

type UserID struct{ PrimaryIDBase }
type UserPreferencesID struct{ PrimaryIDBase }
type MealID struct{ PrimaryIDBase }
type MealPhotoID struct{ PrimaryIDBase }
type TrainingID struct{ PrimaryIDBase }
type TrainingExerciseID struct{ PrimaryIDBase }
type TrainingSetID struct{ PrimaryIDBase }
type ExerciseID struct{ PrimaryIDBase }
type RoutineID struct{ PrimaryIDBase }
type RoutineExerciseID struct{ PrimaryIDBase }
type WeightID struct{ PrimaryIDBase }
