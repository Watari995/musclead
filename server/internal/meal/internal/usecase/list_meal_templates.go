package mealusecase

import (
	"context"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/pagination"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ListMealTemplatesInput struct {
	UserID valueobject.UserID
	Limit  int
	Offset int
}

type ListMealTemplatesOutput struct {
	MealTemplates []*mealdomain.MealTemplate
	Pagination    pagination.OffsetPaginator
}

type ListMealTemplates struct {
	mealTemplateRepo mealdomain.MealTemplateRepository
}

func (uc *ListMealTemplates) Execute(ctx context.Context, input ListMealTemplatesInput) (*ListMealTemplatesOutput, error) {
	mealTemplates, paginator, err := uc.mealTemplateRepo.FindAllByUserIDWithOffsetPagination(ctx, input.UserID, input.Limit, input.Offset)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &ListMealTemplatesOutput{MealTemplates: mealTemplates, Pagination: paginator}, nil
}

func NewListMealTemplates(mealTemplateRepo mealdomain.MealTemplateRepository) *ListMealTemplates {
	return &ListMealTemplates{mealTemplateRepo: mealTemplateRepo}
}
