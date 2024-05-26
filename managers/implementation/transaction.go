package transactionManager

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/nkarakotova/lim-core/managers"
)

type TransactionManagerImplementation struct {
	manager *manager.Manager
}

func NewTransactionManagerImplementation(manager *manager.Manager) managers.TransactionManager {
	return &TransactionManagerImplementation{manager: manager}
}

func (transactor *TransactionManagerImplementation) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return transactor.manager.Do(ctx, fn)
}
