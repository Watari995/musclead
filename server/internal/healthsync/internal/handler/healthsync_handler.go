package healthsynchandler

import (
	"net/http"

	healthsyncusecase "github.com/Watari995/musclead/internal/healthsync/internal/usecase"
	"github.com/Watari995/musclead/internal/shared/httpx"
)

type HealthSyncHandler struct {
	buildAuthURL *healthsyncusecase.BuildAuthURL
	connect      *healthsyncusecase.ConnectHealthPlanet
	frontendURL  string
}

func New(buildAuthURL *healthsyncusecase.BuildAuthURL, connect *healthsyncusecase.ConnectHealthPlanet, frontendURL string) http.Handler {
	h := &HealthSyncHandler{
		buildAuthURL: buildAuthURL,
		connect:      connect,
		frontendURL:  frontendURL,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /integrations/healthplanet/auth", h.Auth)
	mux.HandleFunc("GET /integrations/healthplanet/callback/{token}", h.Connect)
	return mux
}

func (h *HealthSyncHandler) Auth(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	authURL, err := h.buildAuthURL.Execute(healthsyncusecase.BuildAuthURLInput{
		UserID: userID,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"url": authURL})
}

func (h *HealthSyncHandler) Connect(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	code := r.URL.Query().Get("code")

	if err := h.connect.Execute(r.Context(), healthsyncusecase.ConnectHealthPlanetInput{
		Token: token,
		Code:  code,
	}); err != nil {
		httpx.WriteError(w, err)
		return
	}

	http.Redirect(w, r, h.frontendURL+"/settings/integrations?connected=true", http.StatusFound)
}
