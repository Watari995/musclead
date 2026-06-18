package mealdto

import (
	"time"

	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	shareddto "github.com/Watari995/musclead/internal/shared/dto"
)

type UpsertMealTemplateRequest struct {
	Name          string   `json:"name"`
	MealType      string   `json:"meal_type"`
	Calories      int      `json:"calories"`
	ProteinG      *float64 `json:"protein_g,omitempty"`
	FatG          *float64 `json:"fat_g,omitempty"`
	CarbohydrateG *float64 `json:"carbohydrate_g,omitempty"`
}

type UpsertMealTemplateResponse struct {
	MealTemplateID string `json:"meal_template_id"`
}

type ListMealTemplatesResponse struct {
	MealTemplates []MealTemplateDTO       `json:"meal_templates"`
	Pagination    shareddto.PaginationDTO `json:"pagination"`
}

type ReorderMealTemplatesRequest struct {
	MealTemplateIDs []string `json:"meal_template_ids"`
}

type MealTemplateDTO struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Name          string    `json:"name"`
	DisplayOrder  int       `json:"display_order"`
	MealType      string    `json:"meal_type"`
	Calories      int       `json:"calories"`
	ProteinG      *string   `json:"protein_g,omitempty"`
	FatG          *string   `json:"fat_g,omitempty"`
	CarbohydrateG *string   `json:"carbohydrate_g,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewMealTemplateDTOFromEntity(m *mealdomain.MealTemplate) MealTemplateDTO {
	// nullableなvoをstringに変換する
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
	return MealTemplateDTO{
		ID:            m.ID().Value(),
		UserID:        m.UserID().Value(),
		Name:          m.Name().Value(),
		DisplayOrder:  m.DisplayOrder().Value(),
		MealType:      m.MealType().Value(),
		Calories:      m.Calories().Value(),
		ProteinG:      proteinGStr,
		FatG:          fatGStr,
		CarbohydrateG: carbohydrateGStr,
		CreatedAt:     m.CreatedAt(),
		UpdatedAt:     m.UpdatedAt(),
	}
}
