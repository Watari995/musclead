package foodusecase

import (
	"context"

	fooddomain "github.com/Watari995/musclead/internal/food/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
)

type CreateFoodProductInput struct {
	Barcode       *valueobject.Barcode
	Name          valueobject.String100
	Calories      valueobject.NonNegativeInt
	ProteinG      *valueobject.NonNegativeDecimal
	FatG          *valueobject.NonNegativeDecimal
	CarbohydrateG *valueobject.NonNegativeDecimal
}

type CreateFoodProductOutput struct {
	FoodProductID valueobject.FoodProductID
}

// CreateFoodProduct はユーザーが未登録食品を手動登録する。
// register_source は 'user' で保存する。
type CreateFoodProduct struct {
	foodProductRepo fooddomain.FoodProductRepository
}

func (uc *CreateFoodProduct) Execute(ctx context.Context, input CreateFoodProductInput) (*CreateFoodProductOutput, error) {
	foodProduct := fooddomain.CreateFoodProduct(
		input.Barcode,
		input.Name,
		input.Calories,
		input.ProteinG,
		input.FatG,
		input.CarbohydrateG,
		valueobject.NewFoodProductRegisterSourceFromCode(valueobject.FoodProductRegisterSourceUser),
	)
	if err := uc.foodProductRepo.Create(ctx, foodProduct); err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	return &CreateFoodProductOutput{FoodProductID: foodProduct.ID()}, nil
}

func NewCreateFoodProduct(foodProductRepo fooddomain.FoodProductRepository) *CreateFoodProduct {
	return &CreateFoodProduct{foodProductRepo: foodProductRepo}
}
