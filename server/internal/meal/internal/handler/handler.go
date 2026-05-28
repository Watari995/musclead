// Package mealhandler exposes HTTP handlers for the meal module.
package mealhandler

import (
	"net/http"

	mealusecase "github.com/Watari995/musclead/internal/meal/internal/usecase"
)

// Handler は meal モジュールが受け持つ HTTP ルーティングを保持する。
type Handler struct {
	mux *http.ServeMux
}

// New は meal モジュールの HTTP ハンドラを構築する。
//
// TODO: 各ルートに対応する HTTP ハンドラを実装する。
// 現状はすべて 501 Not Implemented を返すプレースホルダー。
func New(
	listMeals *mealusecase.ListMeals,
	deleteMealByID *mealusecase.DeleteMealByID,
) http.Handler {
	h := &Handler{mux: http.NewServeMux()}

	h.mux.HandleFunc("GET /meals", notImplemented("list meals"))
	h.mux.HandleFunc("GET /meals/{id}", notImplemented("find meal"))
	h.mux.HandleFunc("POST /meals", notImplemented("create meal"))
	h.mux.HandleFunc("PUT /meals/{id}", notImplemented("update meal"))
	_ = listMeals
	_ = deleteMealByID

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func notImplemented(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotImplemented)
		_, _ = w.Write([]byte(`{"error":"not implemented: ` + name + `"}`))
	}
}
