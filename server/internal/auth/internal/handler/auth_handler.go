package authhandler

import (
	"net/http"
	"time"

	authdto "github.com/Watari995/musclead/internal/auth/dto"
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

// Login godoc
//
// @Summary ログイン
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "ログイン情報"
// @Success 200 {object} authdto.AccessTokenResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /auth/login [post]
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
	httpx.WriteJSON(w, http.StatusOK, authdto.AccessTokenResponse{
		AccessToken:          output.AccessToken,
		AccessTokenExpiresAt: output.AccessTokenExpiresAt.Format(time.RFC3339),
	})
}

// Refresh godoc
//
// @Summary 認証情報のリフレッシュ
// @Tags auth
// @Produce json
// @Success 200 {object} authdto.AccessTokenResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshRaw, err := httpx.ReadRefreshCookie(r)
	if err != nil {
		httpx.WriteError(w, myerror.NewUnauthorizedError().SetMessage("invalid refresh token"))
		return
	}
	output, err := h.refresh.Execute(r.Context(), authusecase.RefreshInput{
		RefreshRaw: refreshRaw,
		UserAgent:  r.UserAgent(),
		IPAddress:  r.RemoteAddr,
	})
	if err != nil {
		httpx.ClearRefreshCookie(w)
		httpx.WriteError(w, err)
		return
	}
	httpx.SetRefreshCookie(w, output.RefreshToken, output.RefreshTokenExpiresAt)
	httpx.WriteJSON(w, http.StatusOK, authdto.AccessTokenResponse{
		AccessToken:          output.AccessToken,
		AccessTokenExpiresAt: output.AccessTokenExpiresAt.Format(time.RFC3339),
	})
}

// Logout godoc
//
// @Summary ログアウト
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} authdto.AccessTokenResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	refreshRaw, err := httpx.ReadRefreshCookie(r)
	if err != nil {
		httpx.ClearRefreshCookie(w)
		httpx.WriteNoContent(w)
		return
	}
	if err := h.logout.Execute(r.Context(), authusecase.LogoutInput{
		RefreshRaw: refreshRaw,
	}); err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.ClearRefreshCookie(w)
	httpx.WriteNoContent(w)
}
