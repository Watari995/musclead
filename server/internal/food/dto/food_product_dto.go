package fooddto

import (
	"time"

	fooddomain "github.com/Watari995/musclead/internal/food/internal/domain"
)

// FoodProductDTO は食品マスタの API レスポンス。
type FoodProductDTO struct {
	ID             string    `json:"id"`
	Barcode        *string   `json:"barcode,omitempty"`
	Name           string    `json:"name"`
	Calories       int       `json:"calories"`
	ProteinG       *string   `json:"protein_g,omitempty"`
	FatG           *string   `json:"fat_g,omitempty"`
	CarbohydrateG  *string   `json:"carbohydrate_g,omitempty"`
	RegisterSource string    `json:"register_source"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// SearchByNameResponse は食品名検索レスポンス。
type SearchByNameResponse struct {
	FoodProducts []FoodProductDTO `json:"food_products"`
}

// SearchByBarcodeResponse は食品バーコード検索レスポンス。
type SearchByBarcodeResponse struct {
	FoodProducts []FoodProductDTO `json:"food_products"`
}

// CreateFoodProductRequest はユーザーによる食品登録リクエスト。
type CreateFoodProductRequest struct {
	Barcode       *string `json:"barcode,omitempty"`
	Name          string  `json:"name"`
	Calories      int     `json:"calories"`
	ProteinG      *string `json:"protein_g,omitempty"`
	FatG          *string `json:"fat_g,omitempty"`
	CarbohydrateG *string `json:"carbohydrate_g,omitempty"`
}

// CreateFoodProductResponse は食品登録レスポンス。
type CreateFoodProductResponse struct {
	FoodProductID string `json:"food_product_id"`
}

func FromEntity(foodProduct *fooddomain.FoodProduct) FoodProductDTO {
	id := foodProduct.ID().Value()
	name := foodProduct.Name().Value()
	var barcode *string
	if foodProduct.Barcode() != nil {
		s := foodProduct.Barcode().Value()
		barcode = &s
	}
	calories := foodProduct.Calories().Value()
	var proteinG *string
	if foodProduct.ProteinG() != nil {
		s := foodProduct.ProteinG().Value().String()
		proteinG = &s
	}
	var fatG *string
	if foodProduct.FatG() != nil {
		s := foodProduct.FatG().Value().String()
		fatG = &s
	}
	var carbohydrateG *string
	if foodProduct.CarbohydrateG() != nil {
		s := foodProduct.CarbohydrateG().Value().String()
		carbohydrateG = &s
	}
	registerSource := foodProduct.RegisterSource().Value()
	createdAt := foodProduct.CreatedAt()
	updatedAt := foodProduct.UpdatedAt()
	return FoodProductDTO{
		ID:             id,
		Barcode:        barcode,
		Name:           name,
		Calories:       calories,
		ProteinG:       proteinG,
		FatG:           fatG,
		CarbohydrateG:  carbohydrateG,
		RegisterSource: registerSource,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
