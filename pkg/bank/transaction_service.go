package bank

import "bluecoins-to-splitwise-go/pkg/model"

type TransactionService interface {
	ValidFile(filename string) bool
	GetBankTransactions(filename string) ([]model.BankTransaction, error)
	WriteTransactionRecords(records []model.BluecoinsTransactionImport, filename string) error
}
