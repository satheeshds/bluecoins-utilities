package transformer

import (
	"bluecoins-to-splitwise-go/pkg/model"
)

type TransformService interface {
	TrasformBankToBluecoinTransactions(bankTransaction []model.BankTransaction) ([][]string, error)
}
