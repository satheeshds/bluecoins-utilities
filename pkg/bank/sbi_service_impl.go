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

type SBITransactionServiceImpl struct {
	fileService file.FileService
}

func (t *SBITransactionServiceImpl) GetBankTransactions(filename string) ([]model.BankTransaction, error) {
	records, err := t.fileService.ReadContents(filename, '\t')
	if err != nil {
		log.Printf("Error reading transaction records: %v", err)
		return nil, err
	}

	const headerLines = 20
	const trailingLines = 1
	if len(records) <= headerLines {
		log.Printf("No transactions found in file: %s", filename)
		return nil, nil
	}

	// Remove first 20 records as they are headers and a trailing entry
	records = records[headerLines : len(records)-trailingLines]

	transactions := make([]model.BankTransaction, len(records))
	const timeLayout = "2 Jan 2006"
	const timeIndex = 0

	for i, record := range records {
		date, err := time.Parse(timeLayout, record[timeIndex])
		if err != nil {
			log.Printf("Error parsing date (%s): %v", record[timeIndex], err)
			return nil, err
		}

		var transactionType model.TransactionType
		var amount float64
		var amountStr string

		if len(strings.TrimSpace(record[4])) != 0 {
			transactionType = model.Debit
			amountStr = record[4]
		} else {
			transactionType = model.Credit
			amountStr = record[5]
		}

		amountStr = strings.ReplaceAll(amountStr, ",", "")
		amount, err = strconv.ParseFloat(amountStr, 32)

		if err != nil {
			log.Printf("Error parsing amount (%s): %v", amountStr, err)
			return nil, err
		}

		transactions[i] = model.BankTransaction{
			Date:            date,
			Description:     record[2],
			TransactionType: transactionType,
			Amount:          amount,
		}
	}
	return transactions, nil
}

func (t *SBITransactionServiceImpl) WriteTransactionRecords(records []model.BluecoinsTransactionImport, filename string) error {
	return t.fileService.WriteTransactionRecords(records, filename)
}

func (t *SBITransactionServiceImpl) ValidFile(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".xls"
}
