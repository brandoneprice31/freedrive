package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/brandoneprice31/freedrive/config"
	"github.com/brandoneprice31/freedrive/service"
)

type (
	Manager struct {
		config   config.Config
		n        int
		services map[service.ServiceType]service.Service
		bufs     map[service.ServiceType]*Buffer
	}

	services struct {
		services []service.Service
	}

	key struct {
		BackupPath        string              `json:"backup_path"`
		StoredServiceData []storedServiceData `json:"stored_service_data"`
	}

	storedServiceData struct {
		ServiceType service.ServiceType `json:"service_type"`
		ServiceData []byte              `json:"service_data"`
	}

	directory struct {
		directories []directory
		files       []file
		Path        string `json:"directory_path"`
	}

	file struct {
		Data []byte `json:"file_data"`
		size int
		Path string `json:"file_path"`
	}
)

var (
	ErrNoKey = errors.New("no key")
)

func NewManager(c config.Config, ss ...service.Service) *Manager {
	ssMap := make(map[service.ServiceType]service.Service)
	for i := range ss {
		ssMap[ss[i].Type()] = ss[i]
	}
	return &Manager{
		config:   c,
		n:        len(ss),
		services: ssMap,
		bufs:     make(map[service.ServiceType]*Buffer),
	}
}

func (m *Manager) addBuffer(b *Buffer) {
	m.bufs[b.ServiceTye()] = b
}

func loadKey(path string) (*key, error) {
	k := key{}
	data, err := ioutil.ReadFile(path)
	if pErr, ok := err.(*os.PathError); ok && pErr.Err.Error() == "no such file or directory" {
		return &k, ErrNoKey
	} else if err != nil {
		return &k, err
	}

	err = json.Unmarshal(data, &k)
	if err != nil {
		fmt.Printf("err loading key: %s", err.Error())
	}

	return &k, nil
}
