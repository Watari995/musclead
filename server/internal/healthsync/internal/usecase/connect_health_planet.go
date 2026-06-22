package healthsyncusecase

import (
	"context"

	healthsyncdomain "github.com/Watari995/musclead/internal/healthsync/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ConnectHealthPlanetInput struct {
	UserID valueobject.UserID
	// HealthPlanet の OAuth callback で受け取る認可コード
	Code string
}

type ConnectHealthPlanet struct {
	tokenRepo      healthsyncdomain.TokenRepository
	tokenExchanger healthsyncdomain.TokenExchanger
}

func NewConnectHealthPlanet(
	tokenRepo healthsyncdomain.TokenRepository,
	tokenExchanger healthsyncdomain.TokenExchanger,
) *ConnectHealthPlanet {
	return &ConnectHealthPlanet{
		tokenRepo:      tokenRepo,
		tokenExchanger: tokenExchanger,
	}
}

func (uc *ConnectHealthPlanet) Execute(ctx context.Context, input ConnectHealthPlanetInput) error {
	accessToken, refreshToken, expiresAt, err := uc.tokenExchanger.ExchangeCode(ctx, input.Code)
	if err != nil {
		return err
	}

	existing, err := uc.tokenRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return myerror.NewInternalError().Wrap(err)
	}

	var token *healthsyncdomain.Token
	if existing != nil {
		token = existing
		token.UpdateTokens(accessToken, refreshToken, expiresAt)
	} else {
		token = healthsyncdomain.CreateToken(
			input.UserID,
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
