// Package user is the public facade of the user module.
// Modular Monolith (strict) のため、 外部からは Module 経由でのみアクセス可能。
package user

import (
	"net/http"

	shareddomain "github.com/Watari995/musclead/internal/shared/domain"
	"github.com/Watari995/musclead/internal/user/interface/publicfunctions"
	userhandler "github.com/Watari995/musclead/internal/user/internal/handler"
	userinfra "github.com/Watari995/musclead/internal/user/internal/infra"
	userusecase "github.com/Watari995/musclead/internal/user/internal/usecase"
	"github.com/go-gorp/gorp/v3"
)

// Module は user モジュールの公開 API。
// HTTP ハンドラだけを外に出すことで、 内部の usecase / repository を隠蔽する。
type Module struct {
	PublicHandler http.Handler
	Handler       http.Handler
	userCommand   publicfunctions.UserCommand
	userQuery     publicfunctions.UserQuery
}

func NewModule(dbmap *gorp.DbMap, storageClient shareddomain.StorageClient, urlBuilder shareddomain.URLBuilder) *Module {
	// repositoryを作成
	dbmap.AddTableWithName(userinfra.UserModel{}, "users").SetKeys(false, "ID")
	dbmap.AddTableWithName(userinfra.UserPreferencesModel{}, "user_preferences").SetKeys(false, "ID")
	repo := userinfra.NewUserRepository(dbmap)
	prefsRepo := userinfra.NewUserPreferencesRepository(dbmap)
	hasher := userinfra.NewBcryptPasswordHasher()

	register := userusecase.NewRegisterUser(repo, hasher)
	find := userusecase.NewFindUser(repo)
	updateUser := userusecase.NewUpdateUser(repo, storageClient)
	delete := userusecase.NewDeleteUser(repo)
	generateProfileImagePresignedURL := userusecase.NewGenerateProfileImagePresignedURL(storageClient)
	me := userusecase.NewMe(repo, prefsRepo)
	updatePreferences := userusecase.NewUpdatePreferences(prefsRepo)

	authenticate := userusecase.NewAuthenticate(repo, hasher)

	getEmailByUserID := userusecase.NewGetEmailByUserID(repo)

	// 認証済みルートをひとつの mux にまとめる
	authedMux := http.NewServeMux()
	userhandler.RegisterAuthenticatedHandlers(authedMux, urlBuilder, me, find, updateUser, delete, generateProfileImagePresignedURL)
	userhandler.RegisterAuthenticatedPreferencesHandlers(authedMux, updatePreferences)

	return &Module{
		PublicHandler: userhandler.NewPublic(register),
		Handler:       authedMux,
		userCommand:   authenticate,
		userQuery:     getEmailByUserID,
	}
}

// UserCommand は他 module 公開用 getter。
// userCommand は unexported にし setter を持たないことで、 NewModule 後の依存差し替えを防ぐ。
func (m *Module) UserCommand() publicfunctions.UserCommand {
	return m.userCommand
}

// UserQuery は他 module 公開用 getter (UserCommand と同様の immutable 保護)。
func (m *Module) UserQuery() publicfunctions.UserQuery {
	return m.userQuery
}
