package mealusecase

import (
	"context"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ListMealTemplateInput struct {
	UserID valueobject.UserID
	Limit  int
	Offset int
}

type ListMealTemplateOutput struct {
	MealTemplates []*mealdomain.MealTemplate
	Pagination    pagination.OffsetPaginator
}

type ListMealTemplate struct {
	mealTemplateRepo mealdomain.MealTemplateRepository
}

func (uc *ListMealTemplate) Execute(ctx context.Context, input ListMealTemplateInput) (*ListMealTemplateOutput, error) {
	mealTemplates, paginator, err := uc.mealTemplateRepo.FindAllByUserIDWithOffsetPagination(ctx, input.UserID, input.Limit, input.Offset)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &ListMealTemplateOutput{MealTemplates: mealTemplates, Pagination: paginator}, nil
}

func NewListMealTemplate(mealTemplateRepo mealdomain.MealTemplateRepository) *ListMealTemplate {
	return &ListMealTemplate{mealTemplateRepo: mealTemplateRepo}
}
