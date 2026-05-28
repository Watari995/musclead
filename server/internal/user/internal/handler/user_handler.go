package userhandler

import (
	"net/http"

	userusecase "github.com/Watari995/musclead/internal/user/internal/usecase"
)

type UserHandler struct {
	register *userusecase.RegisterUser
	find     *userusecase.FindUser
	delete   *userusecase.DeleteUser
}

func New(register *userusecase.RegisterUser, find *userusecase.FindUser, delete *userusecase.DeleteUser) http.Handler {
	// ServeHTTP interfaceを満たしている必要がある
	h := &UserHandler{
		register: register,
		find:     find,
		delete:   delete,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", h.Register)
	mux.HandleFunc("GET /users/{id}", h.Find)
	mux.HandleFunc("DELETE /users/{id}", h.Delete)
	return mux
}

// Register
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	
}

// Find
func (h *UserHandler) Find(w http.ResponseWriter, r *http.Request) {}

// Delete
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {}
