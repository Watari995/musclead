package foodhandler

import "net/http"

// FoodHandler は食品マスタの HTTP ハンドラ。
//
// Routes:
//
//	GET  /food_products?q={name}       — 名前検索
//	GET  /food_products/barcode/{code} — バーコード検索
//	POST /food_products                — ユーザー登録
type FoodHandler struct {
	// TODO: implement
}

func New() http.Handler {
	// TODO: implement
	panic("not implemented")
}
