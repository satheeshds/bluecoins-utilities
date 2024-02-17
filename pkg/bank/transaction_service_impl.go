package bank

import (
	"bluecoins-to-splitwise-go/pkg/file"
	"bluecoins-to-splitwise-go/pkg/model"
	"log"
	"strconv"
	"strings"
	"time"
)

type TransactionServiceImpl struct {
	fileService file.FileService
}

const timeLayout = "02 Jan 2006"

func NewTransactionService() (*TransactionServiceImpl, error) {
	fileService, err := file.NewFileService()
	if err != nil {
		log.Printf("Error creating file service: %v", err)
	}

	service := &TransactionServiceImpl{
		fileService: fileService,
	}
	return service, nil
}

func (t *TransactionServiceImpl) GetBankTransactions(filename string) ([]model.BankTransaction, error) {
	records, err := t.fileService.ReadTransactionRecords(filename)
	if err != nil {
		log.Printf("Error reading transaction records: %v", err)
		return nil, err
	}

	if len(records) <= 20 {
		log.Printf("No transactions found in file: %s", filename)
		return nil, nil
	}

	// Remove first 20 records as they are headers and a trailing entry
	records = records[20 : len(records)-1]

	transactions := make([]model.BankTransaction, len(records))

	for i, record := range records {
		date, err := time.Parse(timeLayout, record[0])
		if err != nil {
			log.Printf("Error parsing date (%s): %v", record[0], err)
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
			Amount:          float32(amount),
		}
	}
	return transactions, nil
}
