package userhandler

import (
	"net/http"

	userusecase "github.com/Watari995/musclead/internal/user/internal/usecase"
)

type PreferencesHandler struct {
	updatePreferences *userusecase.UpdatePreferences
}

// UpdatePreferences godoc
//
// @Summary プリファレンス更新
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userdto.UpdatePreferencesRequest true "request"
// @Success 200 {object} userdto.UpdatePreferencesResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /users/me/preferences [patch]
func (h *PreferencesHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
}

func NewPreferencesHandler(updatePreferences *userusecase.UpdatePreferences) *PreferencesHandler {
	return &PreferencesHandler{updatePreferences: updatePreferences}
}
