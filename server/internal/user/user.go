// Package user is the public facade of the user module.
// Modular Monolith (strict) のため、 外部からは Module 経由でのみアクセス可能。
package user

import (
	"database/sql"
	"net/http"

	"github.com/Watari995/musclead/internal/user/interface/publicfunctions"
	userhandler "github.com/Watari995/musclead/internal/user/internal/handler"
	userinfra "github.com/Watari995/musclead/internal/user/internal/infra"
	userusecase "github.com/Watari995/musclead/internal/user/internal/usecase"
)

// Module は user モジュールの公開 API。
// HTTP ハンドラだけを外に出すことで、 内部の usecase / repository を隠蔽する。
type Module struct {
	PublicHandler http.Handler
	Handler       http.Handler
	userCommand   publicfunctions.UserCommand
}

func NewModule(db *sql.DB) *Module {
	repo := userinfra.NewUserRepository(db)
	hasher := userinfra.NewBcryptPasswordHasher()

	register := userusecase.NewRegisterUser(repo, hasher)
	find := userusecase.NewFindUser(repo)
	deleteUser := userusecase.NewDeleteUser(repo)

	authenticate := userusecase.NewAuthenticate(repo, hasher)

	return &Module{
		PublicHandler: userhandler.NewPublic(register),
		Handler:       userhandler.NewAuthenticated(find, deleteUser),
		userCommand:   authenticate,
	}
}

// immutableにするためにゲッター経由で取得する
func (m *Module) UserCommand() publicfunctions.UserCommand {
	return m.userCommand
}
