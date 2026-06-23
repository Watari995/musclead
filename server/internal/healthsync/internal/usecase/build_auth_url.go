package healthsyncusecase

import (
	healthsyncdomain "github.com/Watari995/musclead/internal/healthsync/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

const startBase = "https://api.musclead.com/integrations/healthplanet/start"

type BuildAuthURLInput struct {
	UserID valueobject.UserID
}

type BuildAuthURL struct {
	stateSigner healthsyncdomain.StateSigner
}

func NewBuildAuthURL(stateSigner healthsyncdomain.StateSigner) *BuildAuthURL {
	return &BuildAuthURL{stateSigner: stateSigner}
}

func (uc *BuildAuthURL) Execute(input BuildAuthURLInput) (string, error) {
	token, err := uc.stateSigner.Sign(input.UserID)
	if err != nil {
		return "", err
	}
	return startBase + "?token=" + token, nil
}
