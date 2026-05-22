package valueobject

type PrimaryId interface {
	// PrimaryIdBase を埋め込んだ型を指定
	~struct{ PrimaryIdBase }
}

func NewPrimaryId[T PrimaryId]() T {
	return T{PrimaryIdBase: newPrimaryIdBase()}
}

func NewPrimaryIdFromString[T PrimaryId](s string) (*T, error) {
	p, err := newPrimaryIdBaseFromString(s)
	if err != nil {
		return nil, err
	}
	return &T{PrimaryIdBase: p}, nil
}

// ID 型一覧(Shared Kernel として全モジュール共通の参照子)
type UserID struct{ PrimaryIdBase }
type MealID struct{ PrimaryIdBase }
type MealPhotoID struct{ PrimaryIdBase }
type TrainingID struct{ PrimaryIdBase }
type ExerciseID struct{ PrimaryIdBase }
type SetID struct{ PrimaryIdBase }
type RoutineID struct{ PrimaryIdBase }
type WeightID struct{ PrimaryIdBase }
