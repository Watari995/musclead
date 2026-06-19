// Package food は食品マスタモジュールの公開インターフェース。
// バーコード検索・名前検索・ユーザー登録を提供する。
package food

import (
	"net/http"
)

// Module は food モジュールの公開 API。
type Module struct {
	Handler http.Handler
}

func NewModule() *Module {
	// TODO: 依存を注入して handler を組み立てる
	panic("not implemented")
}
