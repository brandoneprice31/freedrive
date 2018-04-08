package service

type (
	Service interface {
		Type() ServiceType
		BufferSize() int
		MaxThreads() int

		NewBackup() error
		Upload([]byte) (*ServiceData, error)

		NewDownload() error
		Download(ServiceData) ([]byte, error)
	}

	ServiceType int

	ServiceData struct {
		Data []byte
	}
)

const (
	localFileSystemServiceType ServiceType = iota
	braintreeServiceType       ServiceType = iota
)
