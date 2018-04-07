package manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"

	"github.com/brandoneprice31/freedrive/service"
)

func (m *manager) Backup(backupPath string) {
	fmt.Printf("starting backup on %s\n", backupPath)

	for _, s := range m.services {
		err := s.NewBackup()
		if err != nil {
			panic(err)
		}
	}

	d, err := loadDirectory(backupPath)
	if err != nil {
		panic(err)
	}

	err = m.uploadContents(d)
	if err != nil {
		panic(err)
	}

	serviceData, err := m.flushBackup()
	if err != nil {
		panic(err)
	}

	err = m.saveServiceData(m.config.KeyPath, serviceData)
	if err != nil {
		panic(err)
	}

	fmt.Println("finished backup")
}

func loadDirectory(path string) (*directory, error) {
	ff, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	d := &directory{
		Path: path,
	}
	for _, f := range ff {
		fpath := fmt.Sprintf("%s/%s", path, f.Name())

		if f.IsDir() {
			// Load this sub directory.
			fd, err := loadDirectory(fpath)
			if err != nil {
				return nil, err
			}
			d.directories = append(d.directories, *fd)

		} else {
			// Attempt to read the file.
			data, err := ioutil.ReadFile(fpath)

			if err != nil && strings.HasSuffix(err.Error(), "is a directory") {
				// You were tricked, this is a sub directory, let's load it.
				fd, lErr := loadDirectory(fpath)
				if lErr != nil {
					return nil, lErr
				}
				d.directories = append(d.directories, *fd)

			} else if err != nil {
				return nil, err

			} else {
				// Append the files.
				d.files = append(d.files, file{
					Data: data,
					size: int(f.Size()),
					Path: fpath,
				})
			}
		}
	}

	return d, nil
}

func (m *manager) uploadContents(d *directory) error {
	return m.uploadDirectory(d, true)
}

func (m *manager) uploadDirectory(d *directory, skip bool) error {
	dirData, err := json.Marshal(d)
	if err != nil {
		return err
	}

	if !skip {
		err = m.uploadData(dirData, len(dirData))
		if err != nil {
			return err
		}
	}

	for _, subD := range d.directories {
		err := m.uploadDirectory(&subD, false)
		if err != nil {
			return err
		}
	}

	for _, f := range d.files {
		fData, err := json.Marshal(f)
		if err != nil {
			return err
		}

		err = m.uploadData(fData, len(fData))
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *manager) flushBackup() ([]storedServiceData, error) {
	var ssd []storedServiceData
	for _, s := range m.services {
		sds, err := s.FlushBackup()
		if err != nil {
			return nil, err
		}

		for _, sd := range sds {
			ssd = append(ssd, storedServiceData{
				ServiceType: s.Type(),
				ServiceData: sd.Data,
			})
		}
	}

	return ssd, nil
}

func (m *manager) saveServiceData(path string, data []storedServiceData) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, jsonData, 0777)
}

func (m *manager) uploadData(data []byte, n int) error {
	return m.randomService().Upload(data)
}

func (m *manager) randomService() service.Service {
	n := len(m.services)
	randn := rand.Intn(n)

	i := 0
	for _, s := range m.services {
		if i == randn {
			return s
		}
		i++
	}

	return m.services[0]
}
