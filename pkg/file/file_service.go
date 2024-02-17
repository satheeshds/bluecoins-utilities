package file

type FileService interface {
	ReadTransactionRecords(filename string) ([][]string, error)
	WriteTransactionRecords(records [][]string, filename string) error
}
