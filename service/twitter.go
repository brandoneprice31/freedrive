package service

import (
	"strconv"

	twitter "github.com/ChimeraCoder/anaconda"
)

type (
	twitterService struct {
		api *twitter.TwitterApi
	}
)

func NewTwitterService(accessToken, accessSecret, consumerKey, consumerSecret string) (Service, error) {
	api := twitter.NewTwitterApiWithCredentials(accessToken, accessSecret, consumerKey, consumerSecret)

	return &twitterService{
		api: api,
	}, nil
}

func (s *twitterService) Type() ServiceType {
	return twitterServiceType
}

func (s *twitterService) BufferSize() int {
	return 120
}

func (s *twitterService) MaxThreads() int {
	return 5
}

func (s *twitterService) MaxStorageSize() int {
	return s.BufferSize() * 150
}

func (s *twitterService) NewBackup() error {
	return nil
}

func (s *twitterService) Remove(sds []ServiceData) error {
	for _, sd := range sds {
		id, err := strconv.Atoi(string(sd.Data))
		if err != nil {
			return err
		}

		s.api.DeleteTweet(int64(id), true)
	}

	return nil
}

func (s *twitterService) Upload(data []byte) (*ServiceData, error) {
	t, err := s.api.PostTweet(string(data), nil)
	if err != nil {
		return nil, err
	}

	return &ServiceData{
		Data: []byte(strconv.Itoa(int(t.Id))),
	}, nil
}

func (s *twitterService) NewDownload() error {
	return nil
}

func (s *twitterService) Download(sd ServiceData) ([]byte, error) {
	id, err := strconv.Atoi(string(sd.Data))
	if err != nil {
		return nil, err
	}

	t, err := s.api.GetTweet(int64(id), nil)
	if err != nil {
		return nil, err
	}

	return []byte(t.Text), nil
}
