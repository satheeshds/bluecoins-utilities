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

func (b *BluecoinsServiceImpl) GetTransactionsAfter(after time.Time) ([]model.Transaction, error) {
	return b.db.GetTransactions(after)
}
