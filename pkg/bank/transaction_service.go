package bank

import "bluecoins-to-splitwise-go/pkg/model"

type TransactionService interface {
	GetBankTransactions(filename string) ([]model.BankTransaction, error)
	WriteTransactionRecords(records []model.BluecoinsTransactionImport, filename string)
}
