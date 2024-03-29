package model

import (
	"fmt"
	"strings"
)

type BluecoinsTransactionImport struct {
	Type           string   `csv:"(1)Type"`
	Date           string   `csv:"(2)Date"`
	ItemOrPayee    string   `csv:"(3)Item or Payee"`
	Amount         string   `csv:"(4)Amount"`
	ParentCategory string   `csv:"(5)Parent Category"`
	Category       string   `csv:"(6)Category"`
	AccountType    string   `csv:"(7)Account Type"`
	Account        string   `csv:"(8)Account"`
	Notes          string   `csv:"(9)Notes"`
	Labels         []string `csv:"(10) Label"`
	Status         string   `csv:"(11) Status"`
	Split          string   `csv:"(12) Split"`
}

func (t BluecoinsTransactionImport) String() string {
	return fmt.Sprintf("%s|%s|%s|%s", t.ItemOrPayee, t.Category, t.ParentCategory, strings.Join(t.Labels, ","))
}

func (t BluecoinsTransactionImport) ToSlice() []string {
	return []string{
		t.Type,
		t.Date,
		t.ItemOrPayee,
		t.Amount,
		t.ParentCategory,
		t.Category,
		t.AccountType,
		t.Account,
		t.Notes,
		strings.Join(t.Labels, " "),
		t.Status,
		t.Split,
	}
}

func (t BluecoinsTransactionImport) IsTransfer() bool {
	return t.Category == "(Transfer)"
}

func (b *BluecoinsTransactionImport) GetTransferTransactions(account Account) []BluecoinsTransactionImport {
	results := []BluecoinsTransactionImport{}
	txntype := b.Type
	b.Type = "t"
	b.Category = "(Transfer)"
	b.ParentCategory = "(Transfer)"
	transfer := *b
	transfer.Account = account.Name
	transfer.AccountType = account.TypeName
	if txntype == "e" {
		b.Amount = "-" + b.Amount
		results = append(results, *b, transfer)
	} else {
		transfer.Amount = "-" + transfer.Amount
		results = append(results, transfer, *b)
	}
	return results
}
