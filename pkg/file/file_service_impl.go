package file

import (
	"bluecoins-to-splitwise-go/pkg/model"
	"encoding/csv"
	"log"
	"os"
	"reflect"
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

func (f *FileServiceImpl) WriteTransactionRecords(records []model.BluecoinsTransactionImport, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	t := reflect.TypeOf(model.BluecoinsTransactionImport{})
	headers := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("csv")
		if tag != "" {
			headers[i] = tag
		}
	}
	err = writer.Write(headers)
	if err != nil {
		log.Fatal(err)
	}

	// Write the records
	for _, record := range records {
		err := writer.Write(record.ToSlice())
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
