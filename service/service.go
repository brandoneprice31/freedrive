package service

type (
	Service interface {
		Type() ServiceType

		NewBackup() error
		Upload([]byte) error
		FlushBackup() ([]ServiceData, error)

		NewDownload() error
		Download(ServiceData) error
		FlushDownload() ([][]byte, error)
	}

	ServiceType int

	ServiceData struct {
		Data []byte
	}
)

const (
	localFileSystemServiceType ServiceType = iota
	googlePhotosServiceType    ServiceType = iota
)
