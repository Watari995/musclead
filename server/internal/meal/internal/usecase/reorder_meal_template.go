package mealusecase

import (
	"context"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ReorderMealTemplateInput struct {
	UserID     valueobject.UserID
	OrderedIDs []valueobject.MealTemplateID
}

type ReorderMealTemplate struct {
	mealTemplateRepo mealdomain.MealTemplateRepository
	txManager        dbtx.TransactionManager
}

func (uc *ReorderMealTemplate) Execute(ctx context.Context, input ReorderMealTemplateInput) error {
	mealTemplates, err := uc.mealTemplateRepo.FindAllByUserID(ctx, input.UserID)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}

	// 渡された並び順が、ユーザーの保有する種目全件のちょうど1度ずつの並び替えであることを検証する
	byID := make(map[string]*mealdomain.MealTemplate, len(mealTemplates))
	for _, m := range mealTemplates {
		byID[m.ID().Value()] = m
	}
	if len(input.OrderedIDs) != len(byID) {
		return myerror.NewBadRequestError().SetMessage("ordered ids must contain every meal templates exactly once")
	}
	// set型で重複チェック
	seen := make(map[string]struct{}, len(input.OrderedIDs))
	for _, id := range input.OrderedIDs {
		key := id.Value()
		if _, ok := byID[key]; !ok {
			return myerror.NewBadRequestError().SetMessage("ordered ids contain an unknown meal template")
		}
		if _, dup := seen[key]; dup {
			return myerror.NewBadRequestError().SetMessage("ordered ids contain a duplicate meal template")
		}
		seen[key] = struct{}{}
	}

	return uc.txManager.Processing(ctx, func(txCTX context.Context) error {
		for index, id := range input.OrderedIDs {
			displayOrder, err := valueobject.NewNonNegativeInt(index)
			if err != nil {
				return myerror.NewInternalError().Wrap(err)
			}
			mealTemplate := byID[id.Value()]
			mealTemplate.SetDisplayOrder(*displayOrder)
			if err := uc.mealTemplateRepo.Save(txCTX, mealTemplate); err != nil {
				return myerror.NewInternalError().Wrap(err)
			}
		}
		return nil
	})
}

func NewReorderMealTemplate(mealTemplateRepo mealdomain.MealTemplateRepository, txManager dbtx.TransactionManager) *ReorderMealTemplate {
	return &ReorderMealTemplate{mealTemplateRepo: mealTemplateRepo, txManager: txManager}
}
