package file

import "bluecoins-to-splitwise-go/pkg/model"

type FileService interface {
	ReadContents(filename string, separator rune) ([][]string, error)
	WriteTransactionRecords(records []model.BluecoinsTransactionImport, filename string) error
}
