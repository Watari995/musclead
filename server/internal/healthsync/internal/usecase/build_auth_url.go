package healthsyncusecase

import (
	"fmt"
	"net/url"

	healthsyncdomain "github.com/Watari995/musclead/internal/healthsync/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

const (
	healthPlanetAuthURL = "https://www.healthplanet.jp/oauth/auth"
	callbackBase        = "https://api.musclead.com/integrations/healthplanet/callback"
)

type BuildAuthURLInput struct {
	UserID valueobject.UserID
}

type BuildAuthURL struct {
	stateSigner healthsyncdomain.StateSigner
	clientID    string
}

func NewBuildAuthURL(stateSigner healthsyncdomain.StateSigner, clientID string) *BuildAuthURL {
	return &BuildAuthURL{
		stateSigner: stateSigner,
		clientID:    clientID,
	}
}

func (uc *BuildAuthURL) Execute(input BuildAuthURLInput) (string, error) {
	token, err := uc.stateSigner.Sign(input.UserID)
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Set("client_id", uc.clientID)
	params.Set("redirect_uri", callbackBase+"/"+token)
	params.Set("scope", "innerscan")
	params.Set("response_type", "code")
	return fmt.Sprintf("%s?%s", healthPlanetAuthURL, params.Encode()), nil
}
