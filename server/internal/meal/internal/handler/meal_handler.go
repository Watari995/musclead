package mealhandler

import (
	"net/http"

	mealdto "github.com/Watari995/musclead/internal/meal/dto"
	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	mealusecase "github.com/Watari995/musclead/internal/meal/internal/usecase"
	"github.com/Watari995/musclead/internal/myerror"
	shareddomain "github.com/Watari995/musclead/internal/shared/domain"
	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	"github.com/Watari995/musclead/internal/shared/httpx"
	sharedstorage "github.com/Watari995/musclead/internal/shared/infra/storage"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/samber/lo"
)

type MealHandler struct {
	urlBuilder                         shareddomain.URLBuilder
	record                             *mealusecase.RecordMeal
	find                               *mealusecase.FindMealByID
	update                             *mealusecase.UpdateMeal
	delete                             *mealusecase.DeleteMealByID
	list                               *mealusecase.ListMeals
	generateMealPhotoImagePresignedURL *mealusecase.GenerateMealPhotoImagePresignedURL
}

func New(
	urlBuilder shareddomain.URLBuilder,
	record *mealusecase.RecordMeal,
	find *mealusecase.FindMealByID,
	update *mealusecase.UpdateMeal,
	delete *mealusecase.DeleteMealByID,
	list *mealusecase.ListMeals,
	generateMealPhotoImagePresignedURL *mealusecase.GenerateMealPhotoImagePresignedURL,
) http.Handler {
	h := &MealHandler{
		urlBuilder:                         urlBuilder,
		record:                             record,
		find:                               find,
		update:                             update,
		delete:                             delete,
		list:                               list,
		generateMealPhotoImagePresignedURL: generateMealPhotoImagePresignedURL,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /meals", h.Record)
	mux.HandleFunc("GET /meals/{id}", h.Find)
	mux.HandleFunc("PUT /meals/{id}", h.Update)
	mux.HandleFunc("DELETE /meals/{id}", h.Delete)
	mux.HandleFunc("GET /meals", h.List)
	mux.HandleFunc("POST /meals/photos/presigned-url", h.GenerateMealPhotoImagePresignedURL)
	return mux
}

// Record godoc
//
// @Summary 食事記録
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
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
	// 各 photo path が自分のディレクトリ配下か & traversal 含まないか検証
	for _, p := range req.Photos {
		if err := sharedstorage.ValidateUserOwnedImagePath(sharedstorage.ImageKindMeal, userID.Value(), p.ImagePath); err != nil {
			httpx.WriteError(w, err)
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
// @Security BearerAuth
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
	resp := mealdto.NewMealDTO(output.Meal, h.urlBuilder)
	httpx.WriteJSON(w, http.StatusOK, resp)
}

// Update godoc
//
// @Summary 食事更新
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
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
	for _, p := range photos {
		if err := sharedstorage.ValidateUserOwnedImagePath(sharedstorage.ImageKindMeal, userID.Value(), p.ImagePath); err != nil {
			httpx.WriteError(w, err)
			return
		}
	}
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
// @Security BearerAuth
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
// @Security BearerAuth
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
			return mealdto.NewMealDTO(m, h.urlBuilder)
		}),
		Pagination: shareddto.NewPaginationDTO(output.Pagination),
	}
	httpx.WriteJSON(w, http.StatusOK, response)
}

// GenerateMealPhotoImagePresignedURL godoc
//
// @Summary 食事写真のPresigned URL生成
// @Tags meals
// @Security BearerAuth
// @Produce json
// @Param request body mealdto.GenerateMealPhotoImagePresignedURLRequest true "食事写真のPresigned URL生成情報"
// @Success 200 {object} mealdto.GenerateMealPhotoImagePresignedURLResponse "食事写真のPresigned URL生成成功"
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /meals/photos/presigned-url [post]
func (h *MealHandler) GenerateMealPhotoImagePresignedURL(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var req mealdto.GenerateMealPhotoImagePresignedURLRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	contentType, err := valueobject.NewImageContentType(req.ContentType)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid content type"))
		return
	}
	params := mealusecase.GenerateMealPhotoImagePresignedURLInput{
		UserID:      userID,
		ContentType: *contentType,
	}
	output, err := h.generateMealPhotoImagePresignedURL.Execute(r.Context(), params)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, mealdto.GenerateMealPhotoImagePresignedURLResponse{
		URL:  output.URL.Value(),
		Path: output.Path,
	})
}
