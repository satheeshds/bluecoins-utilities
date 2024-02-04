package bluecoins

import "bluecoins-to-splitwise-go/pkg/model"

type BluecoinsService interface {
	// GetTransactions returns all transactions from Bluecoins
	GetTransactions() ([]model.Transaction, error)
	GetTransactionsByDateRange(startDate, endDate string) ([]model.Transaction, error)
}
