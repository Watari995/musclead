package mealusecase

import (
	"context"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
)

type DeleteMealTemplateInput struct {
	MealTemplateID valueobject.MealTemplateID
	UserID         valueobject.UserID
}

type DeleteMealTemplate struct {
	mealTemplateRepo mealdomain.MealTemplateRepository
}

func (uc *DeleteMealTemplate) Execute(ctx context.Context, input DeleteMealTemplateInput) error {
	mealTemplate, err := uc.mealTemplateRepo.FindByIDAndUserID(ctx, input.MealTemplateID, input.UserID)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	if mealTemplate == nil {
		return myerror.NewMealTemplateNotFoundError()
	}
	if err := uc.mealTemplateRepo.DeleteByID(ctx, input.MealTemplateID); err != nil {
		return myerror.NewInternalError().Wrap(err)
	}
	return nil
}

func NewDeleteMealTemplate(mealTemplateRepo mealdomain.MealTemplateRepository) *DeleteMealTemplate {
	return &DeleteMealTemplate{mealTemplateRepo: mealTemplateRepo}
}
