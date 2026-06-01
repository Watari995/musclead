package authhandler

import (
	"net/http"
	"time"

	authusecase "github.com/Watari995/musclead/internal/auth/internal/usecase"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/httpx"
	"github.com/Watari995/musclead/internal/valueobject"
)

type AuthHandler struct {
	login   *authusecase.Login
	refresh *authusecase.Refresh
	logout  *authusecase.Logout
}

func NewAuthHandler(login *authusecase.Login, refresh *authusecase.Refresh, logout *authusecase.Logout) http.Handler {
	h := &AuthHandler{
		login:   login,
		refresh: refresh,
		logout:  logout,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("POST /auth/refresh", h.Refresh)
	mux.HandleFunc("POST /auth/logout", h.Logout)
	return mux
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid email"))
		return
	}
	output, err := h.login.Execute(r.Context(),
		authusecase.LoginInput{
			Email:     *email,
			Password:  req.Password,
			UserAgent: r.UserAgent(), // http.Requestから取得
			IPAddress: r.RemoteAddr,  // http.Requestから取得
		},
	)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	// cookieにrefresh tokenをセットする
	httpx.SetRefreshCookie(w, output.RefreshToken, output.RefreshTokenExpiresAt)
	// access tokenを返す
	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"access_token":            output.AccessToken,
		"access_token_expires_at": output.AccessTokenExpiresAt.Format(time.RFC3339),
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
}
