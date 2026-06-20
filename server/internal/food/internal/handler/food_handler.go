package foodhandler

import (
	"net/http"

	fooddto "github.com/Watari995/musclead/internal/food/dto"
	fooddomain "github.com/Watari995/musclead/internal/food/internal/domain"
	foodusecase "github.com/Watari995/musclead/internal/food/internal/usecase"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/httpx"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/samber/lo"
)

// FoodHandler は食品マスタの HTTP ハンドラ。
//
// Routes:
//
//	GET  /food_products?q={name}       — 名前検索
//	GET  /food_products/barcode/{code} — バーコード検索
//	POST /food_products                — ユーザー登録
type FoodHandler struct {
	searchByName      *foodusecase.SearchByName
	searchByBarcode   *foodusecase.SearchByBarcode
	createFoodProduct *foodusecase.CreateFoodProduct
}

func New(
	searchByName *foodusecase.SearchByName,
	searchByBarcode *foodusecase.SearchByBarcode,
	createFoodProduct *foodusecase.CreateFoodProduct,
) http.Handler {
	h := &FoodHandler{
		searchByName:      searchByName,
		searchByBarcode:   searchByBarcode,
		createFoodProduct: createFoodProduct,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /food_products", h.SearchByName)
	mux.HandleFunc("GET /food_products/barcode/{code}", h.SearchByBarcode)
	mux.HandleFunc("POST /food_products", h.CreateFoodProduct)
	return mux
}

// SearchByName godoc
//
// @Summary 食品名検索
// @Tags food_products
// @Produce json
// @Security BearerAuth
// @Param q query string true "食品名"
// @Success 200 {object} fooddto.SearchFoodProductsResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /food_products [get]
func (h *FoodHandler) SearchByName(w http.ResponseWriter, r *http.Request) {
	name, err := valueobject.NewString100(r.URL.Query().Get("q"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid name"))
		return
	}
	params := foodusecase.SearchByNameInput{
		Name: *name,
	}
	output, err := h.searchByName.Execute(r.Context(), params)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	resp := fooddto.SearchByNameResponse{
		FoodProducts: lo.Map(output.FoodProducts, func(foodProduct *fooddomain.FoodProduct, _ int) fooddto.FoodProductDTO {
			return fooddto.FromEntity(foodProduct)
		}),
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}

// SearchByBarcode godoc
//
// @Summary 食品バーコード検索
// @Tags food_products
// @Produce json
// @Security BearerAuth
// @Param code path string true "食品バーコード"
// @Success 200 {object} fooddto.SearchByBarcodeResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /food_products/barcode/{code} [get]
func (h *FoodHandler) SearchByBarcode(w http.ResponseWriter, r *http.Request) {
	barcode, err := valueobject.NewBarcode(r.PathValue("code"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid barcode"))
		return
	}
	params := foodusecase.SearchByBarcodeInput{
		Barcode: *barcode,
	}
	output, err := h.searchByBarcode.Execute(r.Context(), params)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	resp := fooddto.SearchByBarcodeResponse{
		FoodProducts: lo.Map(output.FoodProducts, func(foodProduct *fooddomain.FoodProduct, _ int) fooddto.FoodProductDTO {
			return fooddto.FromEntity(foodProduct)
		}),
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}

// CreateFoodProduct godoc
//
// @Summary 食品登録
// @Tags food_products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body fooddto.CreateFoodProductRequest true "食品登録"
// @Success 201 {object} fooddto.CreateFoodProductResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /food_products [post]
func (h *FoodHandler) CreateFoodProduct(w http.ResponseWriter, r *http.Request) {
	var req fooddto.CreateFoodProductRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	var err error
	var barcode *valueobject.Barcode
	if req.Barcode != nil {
		barcode, err = valueobject.NewBarcode(*req.Barcode)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid barcode"))
			return
		}
	}
	name, err := valueobject.NewString100(req.Name)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid name"))
		return
	}
	calories, err := valueobject.NewNonNegativeInt(req.Calories)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid calories"))
		return
	}
	var proteinG *valueobject.NonNegativeDecimal
	if req.ProteinG != nil {
		proteinG, err = valueobject.NewNonNegativeDecimalFromString(*req.ProteinG)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid protein g"))
			return
		}
	}
	var fatG *valueobject.NonNegativeDecimal
	if req.FatG != nil {
		fatG, err = valueobject.NewNonNegativeDecimalFromString(*req.FatG)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid fat g"))
			return
		}
	}
	var carbohydrateG *valueobject.NonNegativeDecimal
	if req.CarbohydrateG != nil {
		carbohydrateG, err = valueobject.NewNonNegativeDecimalFromString(*req.CarbohydrateG)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid carbohydrate g"))
			return
		}
	}
	params := foodusecase.CreateFoodProductInput{
		Barcode:       barcode,
		Name:          *name,
		Calories:      *calories,
		ProteinG:      proteinG,
		FatG:          fatG,
		CarbohydrateG: carbohydrateG,
	}
	output, err := h.createFoodProduct.Execute(r.Context(), params)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	resp := fooddto.CreateFoodProductResponse{
		FoodProductID: output.FoodProductID.Value(),
	}
	httpx.WriteJSON(w, http.StatusCreated, resp)
}
