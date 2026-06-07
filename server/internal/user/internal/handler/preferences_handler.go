package userhandler

import (
	"net/http"

	"github.com/Watari995/musclead/internal/myerror"
	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	"github.com/Watari995/musclead/internal/shared/httpx"
	userdto "github.com/Watari995/musclead/internal/user/dto"
	userusecase "github.com/Watari995/musclead/internal/user/internal/usecase"
	"github.com/Watari995/musclead/internal/valueobject"
)

type PreferencesHandler struct {
	updatePreferences *userusecase.UpdatePreferences
}

func RegisterAuthenticatedPreferencesHandlers(mux *http.ServeMux, updatePreferences *userusecase.UpdatePreferences) {
	h := &PreferencesHandler{updatePreferences: updatePreferences}
	mux.HandleFunc("PATCH /users/me/preferences", h.UpdatePreferences)
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
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var req userdto.UpdatePreferencesRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	var themePatch shareddto.Patch[valueobject.Theme]
	if req.Theme.Set {
		if req.Theme.Null {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("theme is required"))
			return
		}
		themePatch.Set = true
		v, err := valueobject.NewTheme(req.Theme.Value)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid theme"))
			return
		}
		themePatch.Value = *v
	}
	output, err := h.updatePreferences.Execute(r.Context(), userusecase.UpdatePreferencesInput{
		UserID: userID,
		Theme:  themePatch,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, userdto.UpdatePreferencesResponse{UserID: output.UserID.Value()})
}
