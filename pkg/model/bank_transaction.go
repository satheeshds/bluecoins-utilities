package model

import "time"

type TransactionType int

const (
	Debit TransactionType = iota
	Credit
)

type BankTransaction struct {
	Date            time.Time
	Description     string
	TransactionType TransactionType
	Amount          float32
}
