package model

import (
	"strings"
	"time"
)

type TransactionType int

const (
	Debit TransactionType = iota
	Credit
)

type BankTransaction struct {
	Date            time.Time
	Description     string
	TransactionType TransactionType
	Amount          float64
	AccountName     string
	AccountType     string
}

func (t *BankTransaction) CleanDescription() string {
	parts := strings.Split(t.Description, "/")
	return strings.Trim(parts[len(parts)-1], "-")
}

func (tt *TransactionType) String() string {
	switch *tt {
	case Debit:
		return "Debit"
	case Credit:
		return "Credit"
	default:
		return "Unknown"
	}
}
