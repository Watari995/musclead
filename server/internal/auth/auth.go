// Package auth はログイン・リフレッシュ・ログアウト等の認証機能を提供する Module Facade。
package auth

import (
	"net/http"
	"os"

	sessioninfra "github.com/Watari995/musclead/internal/auth/internal/infra"
	authusecase "github.com/Watari995/musclead/internal/auth/internal/usecase"
	"github.com/Watari995/musclead/internal/user/interface/publicfunctions"
	"github.com/go-gorp/gorp/v3"
)

type Module struct {
	Handler http.Handler
}

func NewModule(dbmap *gorp.DbMap, userCommand publicfunctions.UserCommand) *Module {
	// repositoryを作成する
	dbmap.AddTableWithName(sessioninfra.SessionModel{}, "sessions").SetKeys(false, "ID")
	repo := sessioninfra.NewSessionRepository(dbmap)
	tokenSigner := sessioninfra.NewJWTSigner(os.Getenv("JWT_SECRET"))

	_ = authusecase.NewLogin(userCommand, repo, tokenSigner)
	_ = authusecase.NewRefresh(repo, tokenSigner)
	// _ = authusecase.NewLogout(repo)
	// _ = authusecase.NewMe(repo)

	return &Module{}
}
