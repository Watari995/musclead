package authusecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	sessiondomain "github.com/Watari995/musclead/internal/auth/internal/domain"
	"github.com/Watari995/musclead/internal/user/interface/publicfunctions"
	"github.com/Watari995/musclead/internal/valueobject"
)

type LoginInput struct {
	Email     valueobject.Email
	Password  string
	UserAgent string
	IPAddress string
}

type LoginOutput struct {
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string // 生データ(Cookieに載せる用)
	RefreshTokenExpiresAt time.Time
}

type Login struct {
	userCommand publicfunctions.UserCommand
	sessionRepo sessiondomain.SessionRepository
	tokenSigner sessiondomain.TokenSigner
}

func (uc *Login) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	userRes, err := uc.userCommand.Authenticate(ctx, publicfunctions.AuthenticateRequest{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	// sessionを作成する
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	refreshRaw := base64.RawURLEncoding.EncodeToString(b)
	// DBに保存するのはSHA-256ハッシュ(生は返すだけ)
	sum := sha256.Sum256([]byte(refreshRaw))
	refreshHash := hex.EncodeToString(sum[:])
	// refresh tokenは7日間有効
	now := time.Now()
	refreshTokenExpiresAt := now.Add(7 * 24 * time.Hour)
	session := sessiondomain.CreateSession(userRes.UserID, refreshHash, input.UserAgent, input.IPAddress, refreshTokenExpiresAt)
	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return nil, err
	}

	accessTokenExpiresAt := now.Add(15 * time.Minute)
	accessToken, err := uc.tokenSigner.SignAccessToken(userRes.UserID, accessTokenExpiresAt)
	if err != nil {
		return nil, err
	}
	return &LoginOutput{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshRaw,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}

func NewLogin(userCommand publicfunctions.UserCommand, sessionRepo sessiondomain.SessionRepository, tokenSigner sessiondomain.TokenSigner) *Login {
	return &Login{userCommand: userCommand, sessionRepo: sessionRepo, tokenSigner: tokenSigner}
}
