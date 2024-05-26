package bank

import (
	"bluecoins-to-splitwise-go/pkg/file"
	"errors"
	"log"
)

func NewTransactionService(filename string) (TransactionService, error) {
	fileService, err := file.NewFileService()
	if err != nil {
		log.Printf("Error creating file service: %v", err)
	}

	sbiService := &SBITransactionServiceImpl{
		fileService: fileService,
	}

	if sbiService.ValidFile(filename) {
		log.Printf("Using SBI service for file: %s", filename)
		return sbiService, nil
	}

	hdfcService := &HdfcFreedomTransactionServiceImpl{
		fileService: fileService,
	}

	if hdfcService.ValidFile(filename) {
		log.Printf("Using HDFC Freedom service for file: %s", filename)
		return hdfcService, nil
	}

	return nil, errors.New("cannot find corresponding bank service for file")
}
