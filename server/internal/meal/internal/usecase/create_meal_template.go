package mealusecase

import (
	"context"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
)

type CreateMealTemplateInput struct {
	UserID        valueobject.UserID
	Name          valueobject.String100
	MealType      valueobject.String20
	Calories      valueobject.NonNegativeInt
	ProteinG      *valueobject.NonNegativeDecimal
	FatG          *valueobject.NonNegativeDecimal
	CarbohydrateG *valueobject.NonNegativeDecimal
}

type CreateMealTemplateOutput struct {
	MealTemplateID valueobject.MealTemplateID
}

type CreateMealTemplate struct {
	mealTemplateRepo mealdomain.MealTemplateRepository
}

func (uc *CreateMealTemplate) Execute(ctx context.Context, input CreateMealTemplateInput) (*CreateMealTemplateOutput, error) {
	next, err := uc.mealTemplateRepo.NextDisplayOrder(ctx, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	displayOrder, err := valueobject.NewNonNegativeInt(next)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	mealTemplate := mealdomain.CreateMealTemplate(input.UserID, input.Name, *displayOrder, input.MealType, input.Calories, input.ProteinG, input.FatG, input.CarbohydrateG)
	if err := uc.mealTemplateRepo.Save(ctx, mealTemplate); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &CreateMealTemplateOutput{MealTemplateID: mealTemplate.ID()}, nil
}

func NewCreateMealTemplate(mealTemplateRepo mealdomain.MealTemplateRepository) *CreateMealTemplate {
	return &CreateMealTemplate{mealTemplateRepo: mealTemplateRepo}
}
