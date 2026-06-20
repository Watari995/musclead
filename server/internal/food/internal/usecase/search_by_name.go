package foodusecase

import (
	"context"

	fooddomain "github.com/Watari995/musclead/internal/food/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
)

type SearchByNameInput struct {
	Name valueobject.String100
}

type SearchByNameOutput struct {
	FoodProducts []*fooddomain.FoodProduct
}

// SearchByName は食品名の前方一致で自社 DB を検索する。
type SearchByName struct {
	foodProductRepo fooddomain.FoodProductRepository
}

func (uc *SearchByName) Execute(ctx context.Context, input SearchByNameInput) (*SearchByNameOutput, error) {
	foodProducts, err := uc.foodProductRepo.FindAllByName(ctx, input.Name)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &SearchByNameOutput{FoodProducts: foodProducts}, nil
}

func NewSearchByName(foodProductRepo fooddomain.FoodProductRepository) *SearchByName {
	return &SearchByName{foodProductRepo: foodProductRepo}
}
