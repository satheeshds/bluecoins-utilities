package bluecoins

import (
	"bluecoins-to-splitwise-go/pkg/model"
	"time"
)

type BluecoinsService interface {
	// GetTransactions returns all transactions from Bluecoins
	GetTransactionsAfter(after time.Time) ([]model.Transaction, error)
}
