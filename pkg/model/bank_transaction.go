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
}

func (t *BankTransaction) CleanDescription() string {
	parts := strings.Split(t.Description, "/")
	return strings.Trim(parts[len(parts)-1], "-")
}
