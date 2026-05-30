package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	sessiondomain "github.com/Watari995/musclead/internal/auth/internal/domain"
	"github.com/Watari995/musclead/internal/myerror"
	userdomain "github.com/Watari995/musclead/internal/user/internal/domain"
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
	userRepo       userdomain.UserRepository
	passwordHasher userdomain.PasswordHasher
	sessionRepo    sessiondomain.SessionRepository
	tokenSigner    sessiondomain.TokenSigner
}

func (uc *Login) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, myerror.NewUserNotFoundError()
	}
	if err := uc.passwordHasher.Compare(input.Password, &user.PasswordHash()); err != nil {
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
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	session := sessiondomain.CreateSession(user.ID(), refreshHash, input.UserAgent, input.IPAddress, expiresAt)
	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return nil, err
	}

	accessToken, err := uc.tokenSigner.SignAccessToken(user.ID(), time.Now().Add(1*time.Hour))
	if err != nil {
		return nil, err
	}
	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshHash,
	}, nil
}
