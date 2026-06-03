package sqlerr

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

// 参考 : https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html#error_er_foreign_key_check_failed
const mysqlErrForeignKey uint16 = 1451

func IsForeignKeyViolation(err error) bool {
	if mysqlErr, ok := errors.AsType[*mysql.MySQLError](err); ok {
		return mysqlErr.Number == mysqlErrForeignKey
	}
	return false
}
