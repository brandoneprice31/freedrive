package service

type (
	Service interface {
		// Type simply returns the ServiceType of this service.
		Type() ServiceType

		// BufferSize returns the maximum amount of data that your service can store
		// in a single request.
		BufferSize() int

		// The maximum storages this service allows.
		MaxStorageSize() int

		// MaxThreads returns the maximum number of concurrent threads that are able
		// to access your service at once.
		MaxThreads() int

		// NewBackup prepares the service for backing up data.
		NewBackup() error

		// Removes data from the service given an array of ServiceData.
		Remove([]ServiceData) error

		// Upload takes in bytes, stores it on your service, and returns the ServiceData
		// necessary to retrieve these bytes.
		Upload([]byte) (*ServiceData, error)

		// NewDownload prepares the service for downloading data.
		NewDownload() error

		// Download takes in ServiceData, retrieves the store bytes, and returns them.
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
	dropboxServiceType         ServiceType = iota
	twitterServiceType         ServiceType = iota
)
