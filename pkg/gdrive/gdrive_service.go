package gdrive

type GDriveService interface {
	// DownloadFile downloads a file from Google Drive
	DownloadFile(fileName string, destination string) error
}
