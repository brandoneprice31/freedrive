package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type (
	localFileSystem struct {
		path       string
		totalFiles int
	}

	last struct {
		Num int `json:"num"`
	}
)

func NewLocalFileSystemService(path string) (Service, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/last", path))
	var lastNum int
	if pErr, ok := err.(*os.PathError); ok && pErr.Err.Error() == "no such file or directory" {
		lastNum = 0
	} else if err != nil {
		return nil, err
	} else {
		var l last
		err = json.Unmarshal(data, &l)
		if err != nil {
			return nil, err
		}

		lastNum = l.Num
	}

	return &localFileSystem{
		path:       path,
		totalFiles: lastNum,
	}, nil
}

func (s *localFileSystem) Type() ServiceType {
	return localFileSystemServiceType
}

func (s *localFileSystem) NewBackup() error {
	err := os.RemoveAll(s.path)
	if err != nil {
		return err
	}

	err = os.Mkdir(s.path, 0777)
	if err != nil {
		return err
	}

	s.totalFiles = 0
	return nil
}

func (s *localFileSystem) Upload(data []byte) (*ServiceData, error) {
	fpath := fmt.Sprintf("%s/%d", s.path, s.totalFiles)

	err := ioutil.WriteFile(fpath, data, 0777)
	if err != nil {
		return nil, err
	}

	s.totalFiles++
	return &ServiceData{
		Data: []byte(strconv.Itoa(s.totalFiles - 1)),
	}, nil
}

func (s *localFileSystem) NewDownload() error {
	return nil
}

func (s *localFileSystem) Download(sd ServiceData) ([]byte, error) {
	i, err := strconv.Atoi(string(sd.Data))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(fmt.Sprintf("%s/%d", s.path, i))
}
