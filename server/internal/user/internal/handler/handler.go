// Package userhandler exposes HTTP handlers for the user module.
package userhandler

import (
	"net/http"

	userusecase "github.com/Watari995/musclead/internal/user/internal/usecase"
)

// Handler は user モジュールが受け持つ HTTP ルーティングを保持する。
type Handler struct {
	mux *http.ServeMux
}

// New は user モジュールの HTTP ハンドラを構築する。
//
// TODO: 各ルートに対応する HTTP ハンドラを実装する。
// 現状はすべて 501 Not Implemented を返すプレースホルダー。
func New(
	register *userusecase.RegisterUser,
	find *userusecase.FindUser,
	delete *userusecase.DeleteUser,
) http.Handler {
	h := &Handler{mux: http.NewServeMux()}

	h.mux.HandleFunc("POST /users", notImplemented("register user"))
	h.mux.HandleFunc("GET /users/{id}", notImplemented("find user"))
	h.mux.HandleFunc("DELETE /users/{id}", notImplemented("delete user"))

	// 引数の usecase は handler 実装時に使う想定で、 現状は未使用。
	_ = register
	_ = find
	_ = delete

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
