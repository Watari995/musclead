package healthsynchandler

import (
	"log/slog"
	"net/http"
	"net/url"

	healthsyncusecase "github.com/Watari995/musclead/internal/healthsync/internal/usecase"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/httpx"
)

const (
	healthPlanetAuthURL = "https://www.healthplanet.jp/oauth/auth"
	callbackURI         = "https://api.musclead.com/integrations/healthplanet/callback"
)

type HealthSyncHandler struct {
	buildAuthURL *healthsyncusecase.BuildAuthURL
	connect      *healthsyncusecase.ConnectHealthPlanet
	clientID     string
	frontendURL  string
}

func New(buildAuthURL *healthsyncusecase.BuildAuthURL, connect *healthsyncusecase.ConnectHealthPlanet, clientID, frontendURL string) http.Handler {
	h := &HealthSyncHandler{
		buildAuthURL: buildAuthURL,
		connect:      connect,
		clientID:     clientID,
		frontendURL:  frontendURL,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /integrations/healthplanet/auth", h.Auth)
	mux.HandleFunc("GET /integrations/healthplanet/start", h.Start)
	mux.HandleFunc("GET /integrations/healthplanet/callback", h.Connect)
	return mux
}

func (h *HealthSyncHandler) Auth(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	startURL, err := h.buildAuthURL.Execute(healthsyncusecase.BuildAuthURLInput{
		UserID:      userID,
		RedirectURL: r.URL.Query().Get("redirect_url"),
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"url": startURL})
}

func (h *HealthSyncHandler) Start(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("missing token"))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "hp_state",
		Value:    token,
		MaxAge:   300,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
	params := url.Values{}
	params.Set("client_id", h.clientID)
	params.Set("redirect_uri", callbackURI)
	params.Set("scope", "innerscan")
	params.Set("response_type", "code")
	http.Redirect(w, r, healthPlanetAuthURL+"?"+params.Encode(), http.StatusFound)
}

func (h *HealthSyncHandler) Connect(w http.ResponseWriter, r *http.Request) {
	if errParam := r.URL.Query().Get("error"); errParam != "" {
		http.Redirect(w, r, h.frontendURL+"/settings/integrations?error="+errParam, http.StatusFound)
		return
	}
	cookie, err := r.Cookie("hp_state")
	if err != nil {
		http.Redirect(w, r, h.frontendURL+"/settings/integrations?error=session_expired", http.StatusFound)
		return
	}
	code := r.URL.Query().Get("code")
	redirectURL, err := h.connect.Execute(r.Context(), healthsyncusecase.ConnectHealthPlanetInput{
		Token: cookie.Value,
		Code:  code,
	})
	if err != nil {
		slog.Error("healthplanet connect failed", "err", err)
		target := h.frontendURL + "/settings/integrations?error=connection_failed"
		if redirectURL != "" {
			target = redirectURL + "?error=connection_failed"
		}
		http.Redirect(w, r, target, http.StatusFound)
		return
	}
	target := h.frontendURL + "/settings/integrations?connected=true"
	if redirectURL != "" {
		target = redirectURL + "?connected=true"
	}
	http.Redirect(w, r, target, http.StatusFound)
}
