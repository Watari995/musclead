package mealusecase

import (
	"context"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
)

type UpdateMealTemplateInput struct {
	MealTemplateID valueobject.MealTemplateID
	UserID         valueobject.UserID
	Name           valueobject.String100
	MealType       valueobject.String20
	Calories       valueobject.NonNegativeInt
	ProteinG       *valueobject.NonNegativeDecimal
	FatG           *valueobject.NonNegativeDecimal
	CarbohydrateG  *valueobject.NonNegativeDecimal
}

type UpdateMealTemplateOutput struct {
	MealTemplateID valueobject.MealTemplateID
}

type UpdateMealTemplate struct {
	mealTemplateRepo mealdomain.MealTemplateRepository
}

func (uc *UpdateMealTemplate) Execute(ctx context.Context, input UpdateMealTemplateInput) (*UpdateMealTemplateOutput, error) {
	mealTemplate, err := uc.mealTemplateRepo.FindByIDAndUserID(ctx, input.MealTemplateID, input.UserID)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	if mealTemplate == nil {
		return nil, myerror.NewMealTemplateNotFoundError()
	}
	params := mealdomain.UpdateMealTemplateParams{
		Name:          input.Name,
		MealType:      input.MealType,
		Calories:      input.Calories,
		ProteinG:      input.ProteinG,
		FatG:          input.FatG,
		CarbohydrateG: input.CarbohydrateG,
	}
	mealTemplate.Update(params)
	if err := uc.mealTemplateRepo.Save(ctx, mealTemplate); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &UpdateMealTemplateOutput{MealTemplateID: mealTemplate.ID()}, nil
}

func NewUpdateMealTemplate(mealTemplateRepo mealdomain.MealTemplateRepository) *UpdateMealTemplate {
	return &UpdateMealTemplate{mealTemplateRepo: mealTemplateRepo}
}
