package bluecoins

import (
	"bluecoins-to-splitwise-go/pkg/model"
	"time"
)

type BluecoinsService interface {
	// GetTransactions returns all transactions from Bluecoins
	GetTransactionsAfter(after time.Time, accountId int) ([]model.BluecoinsTransaction, error)
	GetAccounts() ([]model.Account, error)
	GetTransactionsImportFormatByDescription(desc string) ([]model.BluecoinsTransactionImport, error)
}
