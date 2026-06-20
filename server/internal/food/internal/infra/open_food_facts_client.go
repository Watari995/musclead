package foodinfra

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"

	fooddomain "github.com/Watari995/musclead/internal/food/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/shopspring/decimal"
)

const openFoodFactsEndpoint = "https://world.openfoodfacts.org/api/v2/product/%s.json"

// OpenFoodFactsClient は Open Food Facts API のクライアント。
// GET https://world.openfoodfacts.org/api/v2/product/{barcode}.json
type OpenFoodFactsClient struct {
	client *http.Client
}

func NewOpenFoodFactsClient(client *http.Client) *OpenFoodFactsClient {
	return &OpenFoodFactsClient{client: client}
}

type openFoodFactsResponse struct {
	Status  int `json:"status"`
	Product struct {
		ProductName   string `json:"product_name"`
		ProductNameJa string `json:"product_name_ja"`
		Nutriments    struct {
			EnergyKcalServing    float64 `json:"energy-kcal_serving"`
			EnergyKcal100g       float64 `json:"energy-kcal_100g"`
			ProteinsServing      float64 `json:"proteins_serving"`
			Proteins100g         float64 `json:"proteins_100g"`
			FatServing           float64 `json:"fat_serving"`
			Fat100g              float64 `json:"fat_100g"`
			CarbohydratesServing float64 `json:"carbohydrates_serving"`
			Carbohydrates100g    float64 `json:"carbohydrates_100g"`
		} `json:"nutriments"`
	} `json:"product"`
}

// FetchByBarcode はバーコードで Open Food Facts を検索する。
// 商品が存在しない場合は nil, nil を返す。
func (c *OpenFoodFactsClient) FetchByBarcode(ctx context.Context, barcode valueobject.Barcode) (*fooddomain.FoodProduct, error) {
	url := fmt.Sprintf(openFoodFactsEndpoint, barcode.Value())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// Open Food Facts の利用規約に従い User-Agent を設定する
	req.Header.Set("User-Agent", "musclead/1.0")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp openFoodFactsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	// status=0 は商品未登録
	if resp.Status == 0 {
		return nil, nil
	}

	return toFoodProductFromResponse(barcode, resp)
}

func toFoodProductFromResponse(barcode valueobject.Barcode, resp openFoodFactsResponse) (*fooddomain.FoodProduct, error) {
	// 商品名: 日本語優先
	productName := resp.Product.ProductNameJa
	if productName == "" {
		productName = resp.Product.ProductName
	}
	name, err := valueobject.NewString100(productName)
	if err != nil {
		return nil, err
	}

	n := resp.Product.Nutriments

	// カロリー: serving 優先、なければ 100g 換算
	kcal := n.EnergyKcalServing
	if kcal == 0 {
		kcal = n.EnergyKcal100g
	}
	calories, err := valueobject.NewNonNegativeInt(int(math.Round(kcal)))
	if err != nil {
		return nil, err
	}

	// PFC: serving 優先、なければ 100g 換算。両方 0 なら nil
	proteinG, err := optionalDecimal(n.ProteinsServing, n.Proteins100g)
	if err != nil {
		return nil, err
	}
	fatG, err := optionalDecimal(n.FatServing, n.Fat100g)
	if err != nil {
		return nil, err
	}
	carbohydrateG, err := optionalDecimal(n.CarbohydratesServing, n.Carbohydrates100g)
	if err != nil {
		return nil, err
	}

	registerSource, err := valueobject.NewFoodProductRegisterSourceFromCode(valueobject.FoodProductRegisterSourceOpenFoodFacts)
	if err != nil {
		return nil, err
	}

	return fooddomain.CreateFoodProduct(&barcode, *name, *calories, proteinG, fatG, carbohydrateG, *registerSource), nil
}

// optionalDecimal は serving 値を優先し、0 なら per100g を使う。両方 0 なら nil を返す。
func optionalDecimal(serving, per100g float64) (*valueobject.NonNegativeDecimal, error) {
	v := serving
	if v == 0 {
		v = per100g
	}
	if v == 0 {
		return nil, nil
	}
	return valueobject.NewNonNegativeDecimal(decimal.NewFromFloat(v))
}
