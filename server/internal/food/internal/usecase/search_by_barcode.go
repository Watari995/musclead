package foodusecase

import (
	"context"

	fooddomain "github.com/Watari995/musclead/internal/food/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
)

type SearchByBarcodeInput struct {
	Barcode valueobject.Barcode
}

type SearchByBarcodeOutput struct {
	FoodProducts []*fooddomain.FoodProduct
}

// SearchByBarcode は自社 DB → Open Food Facts の順でバーコード検索する。
// 自社 DB になければ外部 API を呼び、ヒットした場合は DB にキャッシュして返す。
// どちらもなければ not found エラーを返す。
type SearchByBarcode struct {
	foodProductRepo fooddomain.FoodProductRepository
	foodFactClient  fooddomain.FoodFactClient
}

func (uc *SearchByBarcode) Execute(ctx context.Context, input SearchByBarcodeInput) (*SearchByBarcodeOutput, error) {
	foodProducts, err := uc.foodProductRepo.FindAllByBarcode(ctx, input.Barcode)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	// 自社DBにヒットしたらそのまま返す
	if len(foodProducts) > 0 {
		return &SearchByBarcodeOutput{FoodProducts: foodProducts}, nil
	}

	foodProductFromOpenFoodFacts, err := uc.foodFactClient.FetchByBarcode(ctx, input.Barcode)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}
	// Open Food Facts でヒットしたら DB にキャッシュして返す
	if foodProductFromOpenFoodFacts != nil {
		if err := uc.foodProductRepo.Create(ctx, foodProductFromOpenFoodFacts); err != nil {
			return nil, myerror.NewInternalError().Wrap(err)
		}
		// これを空配列に追加して返す
		foodProducts = append(foodProducts, foodProductFromOpenFoodFacts)
	} else {
		// どちらもなければnot found エラーを返す
		return nil, myerror.NewNotFoundError()
	}
	return &SearchByBarcodeOutput{FoodProducts: foodProducts}, nil
}

func NewSearchByBarcode(foodProductRepo fooddomain.FoodProductRepository, foodFactClient fooddomain.FoodFactClient) *SearchByBarcode {
	return &SearchByBarcode{foodProductRepo: foodProductRepo, foodFactClient: foodFactClient}
}
