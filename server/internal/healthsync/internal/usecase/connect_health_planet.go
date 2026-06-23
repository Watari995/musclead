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

// Execute は OAuth コードを検証しトークンを保存する。
// 戻り値の redirectURL は JWT に埋め込まれた callback 先 URL(空文字の場合はデフォルト使用)。
func (uc *ConnectHealthPlanet) Execute(ctx context.Context, input ConnectHealthPlanetInput) (redirectURL string, err error) {
	userID, redirectURL, err := uc.stateSigner.Verify(input.Token)
	if err != nil {
		return "", myerror.NewUnauthorizedError().SetMessage("invalid state")
	}
	accessToken, refreshToken, expiresAt, err := uc.tokenExchanger.ExchangeCode(ctx, input.Code)
	if err != nil {
		return redirectURL, err
	}

	existing, err := uc.tokenRepo.FindByUserID(ctx, userID)
	if err != nil {
		return redirectURL, myerror.NewInternalError().Wrap(err)
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
		return redirectURL, myerror.NewInternalError().Wrap(err)
	}

	return redirectURL, nil
}
