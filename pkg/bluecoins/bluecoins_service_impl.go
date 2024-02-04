package bluecoins

import (
	"bluecoins-to-splitwise-go/pkg/db"
	"bluecoins-to-splitwise-go/pkg/gdrive"
	"bluecoins-to-splitwise-go/pkg/model"
	"errors"
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

func (b *BluecoinsServiceImpl) GetTransactions() ([]model.Transaction, error) {

	return b.db.GetTransactions(time.Date(2023, time.December, 24, 0, 0, 0, 0, time.UTC))
	// println(string(file))
	// Replace this with your actual implementation
}

func (b *BluecoinsServiceImpl) GetTransactionsByDateRange(startDate, endDate string) ([]model.Transaction, error) {
	// Parse the dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	//print start and end date
	println(start.Day())
	println(end.Day())

	// Replace this with your actual implementation
	// For example, you might query a database for transactions between start and end
	return nil, errors.New("not implemented")
}
