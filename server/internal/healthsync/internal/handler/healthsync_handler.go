package healthsynchandler

import (
	"fmt"
	"net/http"
	"net/url"

	healthsyncusecase "github.com/Watari995/musclead/internal/healthsync/internal/usecase"
	"github.com/Watari995/musclead/internal/shared/httpx"
)

const (
	healthPlanetAuthURL = "https://www.healthplanet.jp/oauth/auth"
	redirectURI         = "https://api.musclead.com/integrations/healthplanet/callback"
)

type HealthSyncHandler struct {
	clientID string
	connect  *healthsyncusecase.ConnectHealthPlanet
}

func New(clientID string, connect *healthsyncusecase.ConnectHealthPlanet) http.Handler {
	h := &HealthSyncHandler{
		clientID: clientID,
		connect:  connect,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /integrations/healthplanet/auth", h.Auth)
	mux.HandleFunc("GET /integrations/healthplanet/callback", h.Connect)
	return mux
}

func (h *HealthSyncHandler) Auth(w http.ResponseWriter, r *http.Request) {
	params := url.Values{}
	params.Set("client_id", h.clientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("scope", "innerscan")
	params.Set("response_type", "code")

	authURL := fmt.Sprintf("%s?%s", healthPlanetAuthURL, params.Encode())
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *HealthSyncHandler) Connect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	if err := h.connect.Execute(r.Context(), healthsyncusecase.ConnectHealthPlanetInput{
		UserID: userID,
		Code:   code,
	}); err != nil {
		httpx.WriteError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, nil)
}
