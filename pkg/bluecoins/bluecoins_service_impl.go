package bluecoins

import (
	"bluecoins-to-splitwise-go/pkg/db"
	"bluecoins-to-splitwise-go/pkg/gdrive"
	"bluecoins-to-splitwise-go/pkg/model"
	"time"
)

type BluecoinsServiceImpl struct {
	gdriveService gdrive.GDriveServiceImpl
	db            db.DBServiceImpl
	// Add any fields you need here, such as a database connection
}

func NewBluecoinsService() (*BluecoinsServiceImpl, error) {

	driveService, err := gdrive.NewGDriveService()
	if err != nil {
		return nil, err
	}
	fileName := "bluecoins.fydb"
	outFile := "bluecoins.fydb"
	err = driveService.DownloadFile(fileName, outFile)
	if err != nil {
		return nil, err
	}

	dbService, err := db.NewDBService(outFile)
	if err != nil {
		return nil, err
	}
	service := &BluecoinsServiceImpl{}
	service.gdriveService = *driveService
	service.db = *dbService
	// Initialize any fields you need here, such as a database connection
	return service, nil
}

func (b *BluecoinsServiceImpl) GetTransactionsAfter(after time.Time, accountId int) ([]model.BluecoinsTransaction, error) {
	return b.db.GetTransactions(after, accountId)
}

func (b *BluecoinsServiceImpl) GetAccounts() ([]model.Account, error) {
	return b.db.GetAccounts()
}

func (b *BluecoinsServiceImpl) GetTransactionsImportFormatByDescription(desc string) ([]model.BluecoinsTransactionImport, error) {
	return b.db.GetTransactionsImportFormatByDescription(desc)
}

func (b *BluecoinsServiceImpl) GetCategories(text string) ([]model.Category, error) {
	return b.db.GetCategories(text)
}
