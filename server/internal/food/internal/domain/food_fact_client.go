package fooddomain

import (
	"context"

	"github.com/Watari995/musclead/internal/valueobject"
)

type FoodFactClient interface {
	FetchByBarcode(ctx context.Context, barcode valueobject.Barcode) (*FoodProduct, error)
}
