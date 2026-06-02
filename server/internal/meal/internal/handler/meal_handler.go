package mealhandler

import (
	"net/http"

	mealdto "github.com/Watari995/musclead/internal/meal/dto"
	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	mealusecase "github.com/Watari995/musclead/internal/meal/internal/usecase"
	"github.com/Watari995/musclead/internal/myerror"
	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	"github.com/Watari995/musclead/internal/shared/httpx"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/samber/lo"
)

type MealHandler struct {
	record     *mealusecase.RecordMeal
	find       *mealusecase.FindMealByID
	update     *mealusecase.UpdateMeal
	delete     *mealusecase.DeleteMealByID
	list       *mealusecase.ListMeals
	cdnBaseURL string
}

func New(
	record *mealusecase.RecordMeal,
	find *mealusecase.FindMealByID,
	update *mealusecase.UpdateMeal,
	delete *mealusecase.DeleteMealByID,
	list *mealusecase.ListMeals,
	cdnBaseURL string,
) http.Handler {
	h := &MealHandler{
		record:     record,
		find:       find,
		update:     update,
		delete:     delete,
		list:       list,
		cdnBaseURL: cdnBaseURL,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /meals", h.Record)
	mux.HandleFunc("GET /meals/{id}", h.Find)
	mux.HandleFunc("PUT /meals/{id}", h.Update)
	mux.HandleFunc("DELETE /meals/{id}", h.Delete)
	mux.HandleFunc("GET /meals", h.List)
	return mux
}

// Record godoc
//
// @Summary 食事記録
// @Tags meals
// @Accept json
// @Produce json
// @Param X-User-ID header string true "リクエスト元 UserID"
// @Param request body mealdto.RecordMealRequest true "食事記録"
// @Success 201 {object} mealdto.RecordMealResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Router /meals [post]
func (h *MealHandler) Record(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var req mealdto.RecordMealRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
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
	var memo *valueobject.String1000
	if req.Memo != nil {
		memo, err = valueobject.NewString1000(*req.Memo)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid memo"))
			return
		}
	}
	photos := lo.Map(req.Photos, func(p mealdto.MealPhotoInput, _ int) mealdomain.PhotoSpec {
		return mealdomain.PhotoSpec{
			ImagePath:    p.ImagePath,
			DisplayOrder: p.DisplayOrder,
		}
	})
	input := mealusecase.RecordMealInput{
		UserID:        userID,
		EatenAt:       req.EatenAt,
		MealType:      *mealType,
		Calories:      *calories,
		ProteinG:      proteinG,
		FatG:          fatG,
		CarbohydrateG: carbohydrateG,
		Memo:          memo,
		Photos:        photos,
	}

	output, err := h.record.Execute(r.Context(), input)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	resp := mealdto.RecordMealResponse{
		MealID: output.MealID.Value(),
	}
	httpx.WriteJSON(w, http.StatusCreated, resp)
}

// Find godoc
//
// @Summary 食事取得
// @Tags meals
// @Produce json
// @Param X-User-ID header string true "リクエスト元 UserID"
// @Param id path string true "対象 MealID"
// @Success 200 {object} mealdto.MealDTO
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Router /meals/{id} [get]
func (h *MealHandler) Find(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	mealID, err := valueobject.NewPrimaryIDFromString[valueobject.MealID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid mealID"))
		return
	}
	input := mealusecase.FindMealByIDInput{
		MealID: *mealID,
		UserID: userID,
	}
	output, err := h.find.Execute(r.Context(), input)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	resp := mealdto.NewMealDTO(output.Meal, h.cdnBaseURL)
	httpx.WriteJSON(w, http.StatusOK, resp)
}

// Update godoc
//
// @Summary 食事更新
// @Tags meals
// @Accept json
// @Produce json
// @Param X-User-ID header string true "リクエスト元 UserID"
// @Param id path string true "対象 MealID"
// @Param request body mealdto.UpdateMealRequest true "更新内容"
// @Success 200 {object} mealdto.UpdateMealResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Router /meals/{id} [put]
func (h *MealHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	mealID, err := valueobject.NewPrimaryIDFromString[valueobject.MealID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid mealID"))
		return
	}
	var req mealdto.UpdateMealRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
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
	var memo *valueobject.String1000
	if req.Memo != nil {
		memo, err = valueobject.NewString1000(*req.Memo)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid memo"))
			return
		}
	}
	photos := lo.Map(req.Photos, func(p mealdto.MealPhotoInput, _ int) mealdomain.PhotoSpec {
		return mealdomain.PhotoSpec{
			ImagePath:    p.ImagePath,
			DisplayOrder: p.DisplayOrder,
		}
	})
	input := mealusecase.UpdateMealInput{
		MealID:        *mealID,
		UserID:        userID,
		EatenAt:       req.EatenAt,
		MealType:      *mealType,
		Calories:      *calories,
		ProteinG:      proteinG,
		FatG:          fatG,
		CarbohydrateG: carbohydrateG,
		Memo:          memo,
		Photos:        photos,
	}
	output, err := h.update.Execute(r.Context(), input)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	resp := mealdto.UpdateMealResponse{
		MealID: output.MealID.Value(),
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}

// Delete godoc
//
// @Summary 食事削除
// @Tags meals
// @Param X-User-ID header string true "リクエスト元 UserID"
// @Param id path string true "対象 MealID"
// @Success 204
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Router /meals/{id} [delete]
func (h *MealHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	mealID, err := valueobject.NewPrimaryIDFromString[valueobject.MealID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid mealID"))
		return
	}
	input := mealusecase.DeleteMealByIDInput{
		MealID: *mealID,
		UserID: userID,
	}
	if err := h.delete.Execute(r.Context(), input); err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteNoContent(w)
}

// List godoc
//
// @Summary 食事一覧
// @Tags meals
// @Produce json
// @Param X-User-ID header string true "リクエスト元 UserID"
// @Param limit query int false "1ページの件数 (default: 20, max: 100)"
// @Param offset query int false "開始位置 (default: 0)"
// @Success 200 {object} mealdto.ListMealsResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /meals [get]
func (h *MealHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
		// return あり
	}
	limit, offset := httpx.ParseOffsetPagination(r)
	input := mealusecase.ListMealsInput{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	}
	output, err := h.list.Execute(r.Context(), input)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	response := mealdto.ListMealsResponse{
		Meals: lo.Map(output.Meals, func(m *mealdomain.Meal, _ int) mealdto.MealDTO {
			return mealdto.NewMealDTO(m, h.cdnBaseURL)
		}),
		Pagination: shareddto.NewPaginationDTO(output.Pagination),
	}
	httpx.WriteJSON(w, http.StatusOK, response)
}
