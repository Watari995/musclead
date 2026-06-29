package mealdto

import (
	"time"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	shareddomain "github.com/Watari995/musclead/internal/shared/domain"
	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	"github.com/samber/lo"
)

// ─── Request / Response (HTTP境界) ─────────────────────────

type MealPhotoInput struct {
	ImagePath    string `json:"image_path"`
	DisplayOrder int    `json:"display_order"`
}

type RecordMealRequest struct {
	EatenAt       time.Time        `json:"eaten_at"`
	MealType      string           `json:"meal_type"`
	Calories      int              `json:"calories"`
	ProteinG      *float64         `json:"protein_g,omitempty"`
	FatG          *float64         `json:"fat_g,omitempty"`
	CarbohydrateG *float64         `json:"carbohydrate_g,omitempty"`
	Memo          *string          `json:"memo,omitempty"`
	FoodProductID *string          `json:"food_product_id,omitempty"`
	ServingCount  *float64         `json:"serving_count,omitempty"`
	Photos        []MealPhotoInput `json:"photos"`
}

type RecordMealResponse struct {
	MealID string `json:"meal_id"`
}

type UpdateMealRequest struct {
	EatenAt       time.Time        `json:"eaten_at"`
	MealType      string           `json:"meal_type"`
	Calories      int              `json:"calories"`
	ProteinG      *float64         `json:"protein_g,omitempty"`
	FatG          *float64         `json:"fat_g,omitempty"`
	CarbohydrateG *float64         `json:"carbohydrate_g,omitempty"`
	Memo          *string          `json:"memo,omitempty"`
	FoodProductID *string          `json:"food_product_id,omitempty"`
	ServingCount  *float64         `json:"serving_count,omitempty"`
	Photos        []MealPhotoInput `json:"photos"`
}

type UpdateMealResponse struct {
	MealID string `json:"meal_id"`
}

type ListMealsResponse struct {
	Meals      []MealDTO               `json:"meals"`
	Pagination shareddto.PaginationDTO `json:"pagination"`
}

type GenerateMealPhotoImagePresignedURLRequest struct {
	ContentType string `json:"content_type"`
}

type GenerateMealPhotoImagePresignedURLResponse struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

// ─── Entity view ────────────────────────────────────────

type PhotoDTO struct {
	ImagePath    string `json:"image_path"`
	ImageURL     string `json:"image_url"`
	DisplayOrder int    `json:"display_order"`
}

func PhotoFromEntity(p mealdomain.PhotoSpec, urlBuilder shareddomain.URLBuilder) PhotoDTO {
	return PhotoDTO{
		ImagePath:    p.ImagePath,
		ImageURL:     urlBuilder.BuildPublicURL(p.ImagePath),
		DisplayOrder: p.DisplayOrder,
	}
}

type MealDTO struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	EatenAt       time.Time  `json:"eaten_at"`
	MealType      string     `json:"meal_type"`
	Calories      int        `json:"calories"`
	ProteinG      *string    `json:"protein_g,omitempty"`
	FatG          *string    `json:"fat_g,omitempty"`
	CarbohydrateG *string    `json:"carbohydrate_g,omitempty"`
	Memo          *string    `json:"memo,omitempty"`
	FoodProductID *string    `json:"food_product_id,omitempty"`
	ServingCount  string     `json:"serving_count"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Photos        []PhotoDTO `json:"photos"`
}

func FromEntity(m *mealdomain.Meal, urlBuilder shareddomain.URLBuilder) MealDTO {
	// nullable な voをstringに変換する
	var proteinGStr *string
	if m.ProteinG() != nil {
		s := m.ProteinG().Value().String()
		proteinGStr = &s
	}
	var fatGStr *string
	if m.FatG() != nil {
		s := m.FatG().Value().String()
		fatGStr = &s
	}
	var carbohydrateGStr *string
	if m.CarbohydrateG() != nil {
		s := m.CarbohydrateG().Value().String()
		carbohydrateGStr = &s
	}
	var memoStr *string
	if m.Memo() != nil {
		s := m.Memo().Value()
		memoStr = &s
	}

	var foodProductIDStr *string
	if m.FoodProductID() != nil {
		s := m.FoodProductID().Value()
		foodProductIDStr = &s
	}

	photos := lo.Map(m.Photos(), func(p mealdomain.PhotoSpec, idx int) PhotoDTO {
		return PhotoFromEntity(p, urlBuilder)
	})

	return MealDTO{
		ID:            m.ID().Value(),
		UserID:        m.UserID().Value(),
		EatenAt:       m.EatenAt(),
		MealType:      m.MealType().Value(),
		Calories:      m.Calories().Value(),
		ProteinG:      proteinGStr,
		FatG:          fatGStr,
		CarbohydrateG: carbohydrateGStr,
		Memo:          memoStr,
		FoodProductID: foodProductIDStr,
		ServingCount:  m.ServingCount().Value().String(),
		CreatedAt:     m.CreatedAt(),
		UpdatedAt:     m.UpdatedAt(),
		Photos:        photos,
	}
}
