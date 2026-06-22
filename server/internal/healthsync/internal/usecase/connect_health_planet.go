package healthsyncusecase

import (
	"context"

	healthsyncdomain "github.com/Watari995/musclead/internal/healthsync/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
)

type ConnectHealthPlanetInput struct {
	Token string
	Code  string
}

type ConnectHealthPlanet struct {
	tokenRepo      healthsyncdomain.TokenRepository
	tokenExchanger healthsyncdomain.TokenExchanger
	stateSigner    healthsyncdomain.StateSigner
}

func NewConnectHealthPlanet(
	tokenRepo healthsyncdomain.TokenRepository,
	tokenExchanger healthsyncdomain.TokenExchanger,
	stateSigner healthsyncdomain.StateSigner,
) *ConnectHealthPlanet {
	return &ConnectHealthPlanet{
		tokenRepo:      tokenRepo,
		tokenExchanger: tokenExchanger,
		stateSigner:    stateSigner,
	}
}

func (uc *ConnectHealthPlanet) Execute(ctx context.Context, input ConnectHealthPlanetInput) error {
	userID, err := uc.stateSigner.Verify(input.Token)
	if err != nil {
		return myerror.NewUnauthorizedError().SetMessage("invalid state")
	}
	redirectURI := "https://api.musclead.com/integrations/healthplanet/callback/" + input.Token
	accessToken, refreshToken, expiresAt, err := uc.tokenExchanger.ExchangeCode(ctx, input.Code, redirectURI)
	if err != nil {
		return err
	}

	existing, err := uc.tokenRepo.FindByUserID(ctx, userID)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}

	var token *healthsyncdomain.Token
	if existing != nil {
		token = existing
		token.UpdateTokens(accessToken, refreshToken, expiresAt)
	} else {
		token = healthsyncdomain.CreateToken(
			userID,
			accessToken,
			refreshToken,
			expiresAt,
		)
	}
	if err := uc.tokenRepo.Save(ctx, token); err != nil {
		return myerror.NewInternalError().Wrap(err)
	}

	return nil
}
