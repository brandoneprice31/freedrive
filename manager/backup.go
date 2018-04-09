package manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
)

func (m *Manager) Backup(backupPath string) {
	fmt.Printf("starting backup on %s\n", backupPath)

	for _, s := range m.services {
		b, err := NewBackupBuffer(s)
		if err != nil {
			panic(err)
		}

		m.addBuffer(b)
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

	err = m.saveServiceData(backupPath, m.config.KeyPath, serviceData)
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
					continue
				}
				d.directories = append(d.directories, *fd)

			} else if err != nil {
				continue
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

func (m *Manager) uploadContents(d *directory) error {
	return m.uploadDirectory(d, true)
}

func (m *Manager) uploadDirectory(d *directory, skip bool) error {
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

func (m *Manager) flushBackup() ([]storedServiceData, error) {
	var ssd []storedServiceData
	for _, b := range m.bufs {
		sds, err := b.FlushBackup()
		if err != nil {
			return nil, err
		}

		for _, sd := range sds {
			ssd = append(ssd, storedServiceData{
				ServiceType: b.ServiceTye(),
				ServiceData: sd.Data,
			})
		}
	}

	return ssd, nil
}

func (m *Manager) saveServiceData(backupPath, keyPath string, data []storedServiceData) error {
	err := os.Remove(keyPath)
	if err != nil {
		return err
	}

	k := key{
		BackupPath:        backupPath,
		StoredServiceData: data,
	}

	jsonData, err := json.Marshal(k)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(keyPath, jsonData, 0777)
}

func (m *Manager) uploadData(data []byte, n int) error {
	return m.randomBuffer().Save(data)
}

func (m *Manager) randomBuffer() *Buffer {
	n := len(m.bufs)
	randn := rand.Intn(n)

	i := 0
	for _, b := range m.bufs {
		if i == randn {
			return b
		}
		i++
	}

	return m.bufs[0]
}
