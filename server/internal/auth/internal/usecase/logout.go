package authusecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	sessiondomain "github.com/Watari995/musclead/internal/auth/internal/domain"
)

type LogoutInput struct {
	RefreshRaw string
}

type Logout struct {
	sessionRepo sessiondomain.SessionRepository
}

func (uc *Logout) Execute(ctx context.Context, input LogoutInput) error {
	sum := sha256.Sum256([]byte(input.RefreshRaw))
	refreshHash := hex.EncodeToString(sum[:])
	session, err := uc.sessionRepo.FindByRefreshHash(ctx, refreshHash)
	if err != nil {
		return err
	}
	// sessionが見つからない場合は何もしない
	if session == nil {
		return nil
	}
	// sessionが有効でない場合は何もしない
	if !session.IsActive() {
		return nil
	}

	session.Revoke()
	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return err
	}
	return nil
}

func NewLogout(sessionRepo sessiondomain.SessionRepository) *Logout {
	return &Logout{sessionRepo: sessionRepo}
}
