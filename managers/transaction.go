package managers

import (
	"context"
)

//go:generate mockgen  -source=transaction.go -destination=mocks/transaction.go
type TransactionManager interface {
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
}