package valueobject

type PrimaryID interface {
	// PrimaryIDBase を埋め込んだ型を指定
	~struct{ PrimaryIDBase }
	// Genericsの型制約は型情報のみで、埋め込みメソッドを表に出さない。
	// 制約を明示するとTからメソッド呼び出しが可能になる
	Bytes() ([]byte, error)
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
type PaymentID struct{ PrimaryIDBase }
type PaymentEventID struct{ PrimaryIDBase }
type StripeEventID struct{ PrimaryIDBase }
type OutboxEventID struct{ PrimaryIDBase }
type SubscriptionOrderID struct{ PrimaryIDBase }
type SubscriptionID struct{ PrimaryIDBase }
type MealTemplateID struct{ PrimaryIDBase }
type FoodProductID struct{ PrimaryIDBase }
type TokenID struct{ PrimaryIDBase }
type UserWeeklyGoalID struct{ PrimaryIDBase }
type NotificationID struct{ PrimaryIDBase }
