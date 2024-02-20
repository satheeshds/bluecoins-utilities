package transformer

import (
	"bluecoins-to-splitwise-go/pkg/model"
	"strconv"
)

type TransformServiceImp struct{}

func (t *TransformServiceImp) TrasformBankToBluecoinTransactions(bankTransaction []model.BankTransaction) ([][]string, error) {
	bluecoinsTransactions := make([][]string, len(bankTransaction)+1)

	//add header to the first row
	bluecoinsTransactions[0] = []string{
		"(1)Type", "(2)Date", "(3)Item or Payee", "(4)Amount", "(5)Parent Category", "(6)Category", "(7)Account Type", "(8)Account", "(9)Notes", "(10) Label", "(11) Status", "(12) Split",
	}
	for i, bankTransaction := range bankTransaction {
		tranType := "e"
		if bankTransaction.TransactionType == model.Credit {
			tranType = "i"
		}
		bluecoinsTransaction := []string{
			tranType,
			bankTransaction.Date.Format("01/02/2006"),
			bankTransaction.Description, // TODO: Add a function to clean the description
			strconv.FormatFloat(bankTransaction.Amount, 'f', 2, 32),
			"",               // Parent Category
			"",               // Category
			"Bank",           // Account Type
			"SBI Technopark", // Account
			"",               // Notes
			"",               // Label
			"",               // Status
			"",               // Split
		}
		bluecoinsTransactions[i+1] = bluecoinsTransaction
	}
	return bluecoinsTransactions, nil
}
