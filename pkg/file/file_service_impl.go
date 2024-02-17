package file

import (
	"encoding/csv"
	"log"
	"os"
)

type FileServiceImpl struct{}

func NewFileService() (*FileServiceImpl, error) {
	service := &FileServiceImpl{}
	return service, nil
}

func (f *FileServiceImpl) ReadTransactionRecords(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	// Set to -1 to allow for variable number of fields per record
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading all records: %v", err)
		log.Fatal(err)
	}

	return records, nil
}

func (f *FileServiceImpl) WriteTransactionRecords(records [][]string, filename string) error {
	return nil
}
