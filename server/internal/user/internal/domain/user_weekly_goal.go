package userdomain

import (
	"time"

	"github.com/Watari995/musclead/internal/valueobject"
)

type UserWeeklyGoal struct {
	id             valueobject.UserWeeklyGoalID
	userID         valueobject.UserID
	trainingCount  *valueobject.NonNegativeInt
	calorieAverage *valueobject.NonNegativeInt
	weightChangeKg *valueobject.WeightChangeKg
	createdAt      time.Time
	updatedAt      time.Time
}

func (u *UserWeeklyGoal) ID() valueobject.UserWeeklyGoalID {
	return u.id
}

func (u *UserWeeklyGoal) UserID() valueobject.UserID {
	return u.userID
}

func (u *UserWeeklyGoal) TrainingCount() *valueobject.NonNegativeInt {
	return u.trainingCount
}

func (u *UserWeeklyGoal) SetTrainingCount(trainingCount *valueobject.NonNegativeInt) {
	u.trainingCount = trainingCount
	u.updatedAt = time.Now()
}

func (u *UserWeeklyGoal) CalorieAverage() *valueobject.NonNegativeInt {
	return u.calorieAverage
}

func (u *UserWeeklyGoal) SetCalorieAverage(calorieAverage *valueobject.NonNegativeInt) {
	u.calorieAverage = calorieAverage
	u.updatedAt = time.Now()
}

func (u *UserWeeklyGoal) WeightChangeKg() *valueobject.WeightChangeKg {
	return u.weightChangeKg
}

func (u *UserWeeklyGoal) SetWeightChangeKg(weightChangeKg *valueobject.WeightChangeKg) {
	u.weightChangeKg = weightChangeKg
	u.updatedAt = time.Now()
}

func (u *UserWeeklyGoal) CreatedAt() time.Time {
	return u.createdAt
}

func (u *UserWeeklyGoal) UpdatedAt() time.Time {
	return u.updatedAt
}

func CreateUserWeeklyGoal(
	userID valueobject.UserID,
	trainingCount *valueobject.NonNegativeInt,
	calorieAverage *valueobject.NonNegativeInt,
	weightChangeKg *valueobject.WeightChangeKg,
) *UserWeeklyGoal {
	now := time.Now()
	return &UserWeeklyGoal{
		id:             valueobject.NewPrimaryID[valueobject.UserWeeklyGoalID](),
		userID:         userID,
		trainingCount:  trainingCount,
		calorieAverage: calorieAverage,
		weightChangeKg: weightChangeKg,
		createdAt:      now,
		updatedAt:      now,
	}
}

func NewUserWeeklyGoal(
	id valueobject.UserWeeklyGoalID,
	userID valueobject.UserID,
	trainingCount *valueobject.NonNegativeInt,
	calorieAverage *valueobject.NonNegativeInt,
	weightChangeKg *valueobject.WeightChangeKg,
	createdAt time.Time,
	updatedAt time.Time,
) *UserWeeklyGoal {
	return &UserWeeklyGoal{
		id:             id,
		userID:         userID,
		trainingCount:  trainingCount,
		calorieAverage: calorieAverage,
		weightChangeKg: weightChangeKg,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}
