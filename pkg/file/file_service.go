package file

import "bluecoins-to-splitwise-go/pkg/model"

type FileService interface {
	ReadTransactionRecords(filename string) ([][]string, error)
	WriteTransactionRecords(records []model.BluecoinsTransactionImport, filename string) error
}
