package userhandler

import (
	"net/http"

	"time"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/httpx"
	userdto "github.com/Watari995/musclead/internal/user/dto"
	userusecase "github.com/Watari995/musclead/internal/user/internal/usecase"
	"github.com/Watari995/musclead/internal/valueobject"
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

type RegisterRequest struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Birthday *string `json:"birthday,omitempty"`
}

type RegisterResponse struct {
	UserID string `json:"user_id"`
}

// Register
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}

	name, err := valueobject.NewString50(req.Name)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid name"))
		return
	}
	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid email"))
		return
	}
	var birthday *time.Time
	if req.Birthday != nil {
		t, err := time.Parse("2006-01-02", *req.Birthday)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid birthday"))
			return
		}
		birthday = &t
	}

	params := userusecase.RegisterUserInput{
		Name:     *name,
		Email:    *email,
		Birthday: birthday,
		Password: req.Password,
	}

	output, err := h.register.Execute(r.Context(), params)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	resp := RegisterResponse{
		UserID: output.UserID.Value(),
	}
	httpx.WriteJSON(w, http.StatusCreated, resp)
}

// Find
func (h *UserHandler) Find(w http.ResponseWriter, r *http.Request) {
	// path parameterからuserIDを取得
	userID, err := valueobject.NewPrimaryIdFromString[valueobject.UserID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid userID"))
		return
	}
	params := userusecase.FindUserInput{
		UserID: *userID,
	}
	output, err := h.find.Execute(r.Context(), params)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	resp := userdto.NewUserDTO(output.UserID, output.Name, output.Email, output.Birthday, output.CreatedAt, output.UpdatedAt)
	httpx.WriteJSON(w, http.StatusOK, resp)
}

// Delete
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := valueobject.NewPrimaryIdFromString[valueobject.UserID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid userID"))
		return
	}
	if err := h.delete.Execute(r.Context(), userusecase.DeleteUserInput{UserID: *userID}); err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteNoContent(w)
}
