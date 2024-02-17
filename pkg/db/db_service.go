package db

import (
	"bluecoins-to-splitwise-go/pkg/model"
	"time"
)

type DBService interface {
	GetTransactionsAfter(after time.Time, accountid int) ([]model.BluecoinsTransaction, error)
	GetAccounts() ([]model.Account, error)
}
