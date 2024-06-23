package bank

import (
	"bluecoins-to-splitwise-go/pkg/file"
	"bluecoins-to-splitwise-go/pkg/model"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type HdfcFreedomTransactionServiceImpl struct {
	fileService file.FileService
}

func (t *HdfcFreedomTransactionServiceImpl) GetBankTransactions(filename string) ([]model.BankTransaction, error) {
	records, err := t.fileService.ReadContents(filename, '~')
	if err != nil {
		log.Printf("Error reading transaction records: %v", err)
		return nil, err
	}

	const headerLines = 19
	const trailingLines = 12
	if len(records) <= headerLines {
		log.Printf("No transactions found in file: %s", filename)
		return nil, nil
	}

	// Remove the headers and trailing lines
	records = records[headerLines : len(records)-trailingLines]

	transactions := make([]model.BankTransaction, len(records))
	const timeLayout = "02/01/2006 15:04:05"
	const timeIndex = 2

	for i, record := range records {
		date, err := time.Parse(timeLayout, record[timeIndex])
		if err != nil {
			log.Printf("Error parsing date (%s): %v", record[timeIndex], err)
			continue
			// return nil, err
		}

		var transactionType model.TransactionType
		var amount float64
		var amountStr string

		// if len(strings.TrimSpace(record[4])) != 0 {
		// 	transactionType = model.Debit
		// 	amountStr = record[4]
		// } else {
		// 	transactionType = model.Credit
		// 	amountStr = record[5]
		// }

		transactionType = model.Debit
		//TODO: Check credit transactions
		amountStr = record[4]

		amountStr = strings.ReplaceAll(amountStr, ",", "")
		amount, err = strconv.ParseFloat(amountStr, 32)

		if err != nil {
			log.Printf("Error parsing amount (%s): %v", amountStr, err)
			return nil, err
		}

		description := record[3]
		description = strings.TrimPrefix(description, "UPI-")
		transactions[i] = model.BankTransaction{
			Date:            date,
			Description:     description,
			TransactionType: transactionType,
			Amount:          amount,
			AccountName:     "Hdfc Freedom",
			AccountType:     "Credit Card",
		}
		log.Printf("Transaction: %v", transactions[i])
	}

	return transactions, nil
}

func (t *HdfcFreedomTransactionServiceImpl) WriteTransactionRecords(records []model.BluecoinsTransactionImport, filename string) error {
	return t.fileService.WriteTransactionRecords(records, filename)
}

func (t *HdfcFreedomTransactionServiceImpl) ValidFile(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".csv"
}
