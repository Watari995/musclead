package authusecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	sessiondomain "github.com/Watari995/musclead/internal/auth/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
)

type RefreshInput struct {
	RefreshRaw string
	UserAgent  string
	IPAddress  string
}

type RefreshOutput struct {
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

type Refresh struct {
	sessionRepo sessiondomain.SessionRepository
	tokenSigner sessiondomain.TokenSigner
}

func (uc *Refresh) Execute(ctx context.Context, input RefreshInput) (RefreshOutput, error) {
	sum := sha256.Sum256([]byte(input.RefreshRaw))
	refreshHash := hex.EncodeToString(sum[:])
	session, err := uc.sessionRepo.FindByRefreshHash(ctx, refreshHash)
	if err != nil {
		return RefreshOutput{}, err
	}
	if session == nil {
		return RefreshOutput{}, myerror.NewInvalidCredentialsError()
	}
	if !session.IsActive() {
		return RefreshOutput{}, myerror.NewInvalidCredentialsError()
	}

	session.Revoke()
	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return RefreshOutput{}, err
	}

	// sessionを作成する
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return RefreshOutput{}, err
	}
	newRefreshRaw := base64.RawURLEncoding.EncodeToString(b)
	newSum := sha256.Sum256([]byte(newRefreshRaw))
	newRefreshHash := hex.EncodeToString(newSum[:])
	now := time.Now()
	refreshTokenExpiresAt := now.Add(7 * 24 * time.Hour)
	newSession := sessiondomain.CreateSession(session.UserID(), newRefreshHash, input.UserAgent, input.IPAddress, refreshTokenExpiresAt)
	if err := uc.sessionRepo.Save(ctx, newSession); err != nil {
		return RefreshOutput{}, err
	}

	// access tokenを作成する
	accessTokenExpiresAt := now.Add(15 * time.Minute)
	accessToken, err := uc.tokenSigner.SignAccessToken(session.UserID(), accessTokenExpiresAt)
	if err != nil {
		return RefreshOutput{}, err
	}

	return RefreshOutput{AccessToken: accessToken, AccessTokenExpiresAt: accessTokenExpiresAt, RefreshToken: newRefreshRaw, RefreshTokenExpiresAt: refreshTokenExpiresAt}, nil
}

func NewRefresh(sessionRepo sessiondomain.SessionRepository, tokenSigner sessiondomain.TokenSigner) *Refresh {
	return &Refresh{sessionRepo: sessionRepo, tokenSigner: tokenSigner}
}
