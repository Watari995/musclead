package dbtx

import (
	"context"
	"fmt"

	"github.com/go-gorp/gorp/v3"
)

type txKey struct{}

type transactionManager struct {
	dbmap *gorp.DbMap
}

func NewTransactionManager(dbmap *gorp.DbMap) TransactionManager {
	return &transactionManager{dbmap: dbmap}
}

func (m *transactionManager) Processing(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := m.dbmap.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
	}()
	ctx = context.WithValue(ctx, txKey{}, tx)
	if err := f(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Querierはctxからgorp.Transactionを取得し、取得できなければfallbackを返す
func Querier(ctx context.Context, fallback gorp.SqlExecutor) gorp.SqlExecutor {
	if tx, ok := ctx.Value(txKey{}).(*gorp.Transaction); ok {
		return tx
	}

	return fallback
}
