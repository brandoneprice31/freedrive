package service

import (
	"bytes"
	"io/ioutil"
	"log"
	"math/rand"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

type (
	dropboxService struct {
		api files.Client
	}
)

func NewDropboxService(accessToken string) (Service, error) {
	log.SetOutput(ioutil.Discard)

	config := dropbox.Config{
		Token: accessToken,
	}
	c := files.New(config)

	return &dropboxService{
		api: c,
	}, nil
}

func (s *dropboxService) Type() ServiceType {
	return dropboxServiceType
}

func (s *dropboxService) BufferSize() int {
	return 10000000
}

func (s *dropboxService) MaxThreads() int {
	return 5
}

func (s *dropboxService) NewBackup() error {
	return nil
}

func (s *dropboxService) Upload(data []byte) (*ServiceData, error) {
	r := bytes.NewReader(data)
	ci := files.NewCommitInfoWithProperties("/freedrive/" + randomFileName())
	ci.Autorename = true
	f, err := s.api.AlphaUpload(ci, r)
	if err != nil {
		return nil, err
	}

	return &ServiceData{
		Data: []byte(f.PathLower),
	}, nil
}

func randomFileName() string {
	n := rand.Intn(20) + 5
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b) + ".txt"
}

func (s *dropboxService) NewDownload() error {
	return nil
}

func (s *dropboxService) Download(sd ServiceData) ([]byte, error) {
	da := files.NewDownloadArg(string(sd.Data))
	_, r, err := s.api.Download(da)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(r)
}
