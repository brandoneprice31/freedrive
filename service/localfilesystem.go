package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

type (
	localFileSystem struct {
		path       string
		totalFiles int

		serviceData []ServiceData

		buf     []obj
		bufSize int

		dataTag int

		downloaded map[int][]obj
	}

	obj struct {
		Data []byte `json:"data"`
		Tag  int    `json:"tag"`
		Iter int    `json:"iter"`
	}

	objs []obj

	last struct {
		Num int `json:"num"`
	}
)

const (
	MaxBufSize = 1000000
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
		downloaded: make(map[int][]obj),
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

func (s *localFileSystem) Upload(data []byte) error {
	err := s.upload(obj{
		Data: data,
		Tag:  s.dataTag,
		Iter: 0,
	})
	if err != nil {
		return err
	}

	s.dataTag++
	return nil
}

func (s *localFileSystem) upload(o obj) error {
	if len(o.Data) == 0 {
		return nil
	}

	if s.bufSize == MaxBufSize {
		err := s.save()
		if err != nil {
			return err
		}

		s.buf = make([]obj, 0)
		s.bufSize = 0
	}

	spaceLeft := MaxBufSize - s.bufSize

	if spaceLeft >= len(o.Data) {
		s.buf = append(s.buf, o)
		s.bufSize += len(o.Data)
		return nil
	}

	err := s.upload(obj{
		Data: o.Data[:spaceLeft],
		Tag:  o.Tag,
		Iter: o.Iter,
	})
	if err != nil {
		return err
	}

	return s.upload(obj{
		Data: o.Data[spaceLeft:],
		Tag:  o.Tag,
		Iter: o.Iter + 1,
	})
}

func (s *localFileSystem) save() error {
	fpath := fmt.Sprintf("%s/%d", s.path, s.totalFiles)
	data, err := json.Marshal(s.buf)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fpath, data, 0777)
	if err != nil {
		return err
	}

	s.serviceData = append(s.serviceData, ServiceData{
		Data: []byte(strconv.Itoa(s.totalFiles)),
	})
	s.totalFiles++
	return nil
}

func (s *localFileSystem) FlushBackup() ([]ServiceData, error) {
	if s.bufSize > 0 {
		err := s.save()
		if err != nil {
			return nil, err
		}
	}

	data, err := json.Marshal(last{Num: s.totalFiles})
	if err != nil {
		return nil, err
	}
	ioutil.WriteFile(fmt.Sprintf("%s/%s", s.path, "last"), data, 0777)

	return s.serviceData, nil
}

func (s *localFileSystem) NewDownload() error {
	return nil
}

func (s *localFileSystem) Download(sd ServiceData) error {
	i, err := strconv.Atoi(string(sd.Data))
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%d", s.path, i))
	if err != nil {
		return err
	}

	var oo []obj
	err = json.Unmarshal(data, &oo)
	if err != nil {
		return err
	}

	for _, o := range oo {
		s.downloaded[o.Tag] = append(s.downloaded[o.Tag], o)
	}

	return nil
}

func (s *localFileSystem) FlushDownload() ([][]byte, error) {
	dd := [][]byte{}
	for _, oo := range s.downloaded {
		sort.Sort(objs(oo))

		d := []byte{}
		for _, o := range oo {
			d = append(d, o.Data...)
		}

		dd = append(dd, d)
	}

	return dd, nil
}

func (oo objs) Len() int {
	return len(oo)
}

func (oo objs) Less(i, j int) bool {
	return oo[i].Iter < oo[j].Iter
}

func (oo objs) Swap(i, j int) {
	oo[i], oo[j] = oo[j], oo[i]
}
