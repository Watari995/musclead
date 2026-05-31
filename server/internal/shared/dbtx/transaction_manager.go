package dbtx

import "context"

type TransactionManager interface {
	Processing(ctx context.Context, f func(ctx context.Context) error) error
}
