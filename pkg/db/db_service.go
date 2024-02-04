package db

import (
	"bluecoins-to-splitwise-go/pkg/model"
	"time"
)

type DBService interface {
	GetTransactionsAfter(after time.Time) ([]model.Transaction, error)
}
