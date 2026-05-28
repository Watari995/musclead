package mealhandler

import (
	"net/http"
	"time"

	mealdto "github.com/Watari995/musclead/internal/meal/dto"
	mealdomain "github.com/Watari995/musclead/internal/meal/internal/domain"
	mealusecase "github.com/Watari995/musclead/internal/meal/internal/usecase"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/httpx"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
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

type RecordMealRequest struct {
	EatenAt       time.Time `json:"eaten_at"`
	MealType      string    `json:"meal_type"`
	Calories      int       `json:"calories"`
	ProteinG      *float64  `json:"protein_g,omitempty"`
	FatG          *float64  `json:"fat_g,omitempty"`
	CarbohydrateG *float64  `json:"carbohydrate_g,omitempty"`
	Memo          *string   `json:"memo,omitempty"`
	Photos        []struct {
		ImagePath    string `json:"image_path"`
		DisplayOrder int    `json:"display_order"`
	} `json:"photos"`
}

type RecordMealResponse struct {
	MealID string `json:"meal_id"`
}

func (h *MealHandler) Record(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var req RecordMealRequest
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
	var proteinG *valueobject.NonNegativeDecimal
	if req.ProteinG != nil {
		proteinG, err = valueobject.NewNonNegativeDecimal(decimal.NewFromFloat(*req.ProteinG))
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid protein g"))
			return
		}
	}
	var fatG *valueobject.NonNegativeDecimal
	if req.FatG != nil {
		fatG, err = valueobject.NewNonNegativeDecimal(decimal.NewFromFloat(*req.FatG))
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid fat g"))
			return
		}
	}
	var carbohydrateG *valueobject.NonNegativeDecimal
	if req.CarbohydrateG != nil {
		carbohydrateG, err = valueobject.NewNonNegativeDecimal(decimal.NewFromFloat(*req.CarbohydrateG))
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid carbohydrate g"))
			return
		}
	}
	var memo *valueobject.String1000
	if req.Memo != nil {
		memo, err = valueobject.NewString1000(*req.Memo)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid memo"))
			return
		}
	}
	photos := lo.Map(req.Photos, func(p struct {
		ImagePath    string `json:"image_path"`
		DisplayOrder int    `json:"display_order"`
	}, _ int) mealdomain.PhotoData {
		return mealdomain.PhotoData{
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
		PhotoData:     photos,
	}

	output, err := h.record.Execute(r.Context(), input)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	resp := RecordMealResponse{
		MealID: output.MealID.Value(),
	}
	httpx.WriteJSON(w, http.StatusCreated, resp)
}

func (h *MealHandler) Find(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	mealID, err := valueobject.NewPrimaryIdFromString[valueobject.MealID](r.PathValue("id"))
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

func (h *MealHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	mealID, err := valueobject.NewPrimaryIdFromString[valueobject.MealID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid mealID"))
		return
	}
	input := 
}

func (h *MealHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	mealID, err := valueobject.NewPrimaryIdFromString[valueobject.MealID](r.PathValue("id"))
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

func (h *MealHandler) List(w http.ResponseWriter, r *http.Request) {}
