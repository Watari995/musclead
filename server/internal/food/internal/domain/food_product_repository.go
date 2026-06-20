package fooddomain

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

type FoodProductRepository interface {
	FindByID(ctx context.Context, id valueobject.FoodProductID) (*FoodProduct, error)
	FindAllByBarcode(ctx context.Context, barcode valueobject.Barcode) ([]*FoodProduct, error)
	FindAllByName(ctx context.Context, name valueobject.String100) ([]*FoodProduct, error)
	Create(ctx context.Context, foodProduct *FoodProduct) error
}
