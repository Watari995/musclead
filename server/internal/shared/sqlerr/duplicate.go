// Package sqlerr は MySQL ドライバから返るエラーを判別するヘルパーを提供する。
// 各 infra Repository が直接 errno を文字列マッチするのを防ぐため、
// 「重複キー」 等の意味的な判定はこのパッケージに集約する。
package sqlerr

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

// mysqlErrDupEntry は MySQL の ER_DUP_ENTRY(UNIQUE 制約違反)のエラー番号。
//
// 参考: https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html#error_er_dup_entry
const mysqlErrDupEntry uint16 = 1062

// IsDuplicateKey は err が MySQL の UNIQUE 制約違反(ER_DUP_ENTRY)かを判定する。
//
// 用途: INSERT / UPDATE が UNIQUE 制約に当たった時に、 ドメインの sentinel error
// (例: ErrExerciseNameDuplicated)へ変換する分岐に使う。
func IsDuplicateKey(err error) bool {
	if mysqlErr, ok := errors.AsType[*mysql.MySQLError](err); ok {
		return mysqlErr.Number == mysqlErrDupEntry
	}
	return false
}
