package mealhandler

import (
	"net/http"

	mealdto "github.com/Watari995/musclead/internal/meal/dto"
	mealusecase "github.com/Watari995/musclead/internal/meal/internal/usecase"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/httpx"
	"github.com/Watari995/musclead/internal/valueobject"
)

// CreateMealTemplate godoc
//
// @Summary 食事テンプレート作成
// @Tags meal_templates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body mealdto.UpsertMealTemplateRequest true "テンプレート作成"
// @Success 201 {object} mealdto.UpsertMealTemplateResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /meal_templates [post]
func (h *MealHandler) CreateMealTemplate(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var req mealdto.UpsertMealTemplateRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	name, err := valueobject.NewString100(req.Name)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid name"))
		return
	}
	mealType, err := valueobject.NewString20(req.MealType)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid meal type"))
		return
	}
	calories, err := valueobject.NewNonNegativeInt(req.Calories)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid calories"))
		return
	}
	proteinG, err := valueobject.ParseOptionalNonNegativeDecimal(req.ProteinG)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid protein g"))
		return
	}
	fatG, err := valueobject.ParseOptionalNonNegativeDecimal(req.FatG)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid fat g"))
		return
	}
	carbohydrateG, err := valueobject.ParseOptionalNonNegativeDecimal(req.CarbohydrateG)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid carbohydrate g"))
		return
	}
	output, err := h.createMealTemplate.Execute(r.Context(), mealusecase.CreateMealTemplateInput{
		UserID:        userID,
		Name:          *name,
		MealType:      *mealType,
		Calories:      *calories,
		ProteinG:      proteinG,
		FatG:          fatG,
		CarbohydrateG: carbohydrateG,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, mealdto.UpsertMealTemplateResponse{
		MealTemplateID: output.MealTemplateID.Value(),
	})
}

// UpdateMealTemplate godoc
//
// @Summary 食事テンプレート更新
// @Tags meal_templates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "対象 MealTemplateID"
// @Param request body mealdto.UpsertMealTemplateRequest true "更新内容"
// @Success 200 {object} mealdto.UpsertMealTemplateResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Router /meal_templates/{id} [put]
func (h *MealHandler) UpdateMealTemplate(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	templateID, err := valueobject.NewPrimaryIDFromString[valueobject.MealTemplateID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid meal template id"))
		return
	}
	var req mealdto.UpsertMealTemplateRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	name, err := valueobject.NewString100(req.Name)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid name"))
		return
	}
	mealType, err := valueobject.NewString20(req.MealType)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid meal type"))
		return
	}
	calories, err := valueobject.NewNonNegativeInt(req.Calories)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid calories"))
		return
	}
	proteinG, err := valueobject.ParseOptionalNonNegativeDecimal(req.ProteinG)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid protein g"))
		return
	}
	fatG, err := valueobject.ParseOptionalNonNegativeDecimal(req.FatG)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid fat g"))
		return
	}
	carbohydrateG, err := valueobject.ParseOptionalNonNegativeDecimal(req.CarbohydrateG)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid carbohydrate g"))
		return
	}
	output, err := h.updateMealTemplate.Execute(r.Context(), mealusecase.UpdateMealTemplateInput{
		MealTemplateID: *templateID,
		UserID:         userID,
		Name:           *name,
		MealType:       *mealType,
		Calories:       *calories,
		ProteinG:       proteinG,
		FatG:           fatG,
		CarbohydrateG:  carbohydrateG,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, mealdto.UpsertMealTemplateResponse{
		MealTemplateID: output.MealTemplateID.Value(),
	})
}

// DeleteMealTemplate godoc
//
// @Summary 食事テンプレート削除
// @Tags meal_templates
// @Security BearerAuth
// @Param id path string true "対象 MealTemplateID"
// @Success 204
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Router /meal_templates/{id} [delete]
func (h *MealHandler) DeleteMealTemplate(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	templateID, err := valueobject.NewPrimaryIDFromString[valueobject.MealTemplateID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid meal template id"))
		return
	}
	if err := h.deleteMealTemplate.Execute(r.Context(), mealusecase.DeleteMealTemplateInput{
		MealTemplateID: *templateID,
		UserID:         userID,
	}); err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteNoContent(w)
}

// ReorderMealTemplates godoc
//
// @Summary 食事テンプレート並び替え
// @Tags meal_templates
// @Accept json
// @Security BearerAuth
// @Param request body mealdto.ReorderMealTemplatesRequest true "並び替え後のID順"
// @Success 204
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /meal_templates/reorder [post]
func (h *MealHandler) ReorderMealTemplates(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var req mealdto.ReorderMealTemplatesRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	orderedIDs := make([]valueobject.MealTemplateID, 0, len(req.MealTemplateIDs))
	for _, rawID := range req.MealTemplateIDs {
		id, err := valueobject.NewPrimaryIDFromString[valueobject.MealTemplateID](rawID)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid meal template id in list"))
			return
		}
		orderedIDs = append(orderedIDs, *id)
	}
	if err := h.reorderMealTemplates.Execute(r.Context(), mealusecase.ReorderMealTemplateInput{
		UserID:     userID,
		OrderedIDs: orderedIDs,
	}); err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteNoContent(w)
}
