package db

import (
	"bluecoins-to-splitwise-go/pkg/model"
	"time"
)

type DBService interface {
	GetTransactionsAfter(after time.Time, accountid int) ([]model.BluecoinsTransaction, error)
	GetAccounts() ([]model.Account, error)
	GetAccountsBySearch(prefix string) ([]model.Account, error)
	GetTransactionsImportFormatByDescription(desc string) ([]model.BluecoinsTransactionImport, error)
	GetCategories(text string) ([]model.Category, error)
}
