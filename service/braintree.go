package service

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/PuerkitoBio/rehttp"
	"github.com/lionelbarrow/braintree-go"
)

type (
	braintreeService struct {
		account account
		api     *braintree.Braintree
	}

	account struct {
		merchID string
		pubKey  string
		privKey string
	}
)

const (
	MaxTries        = 5
	MaxTimeout      = 5 * time.Second
	NumCustomFields = 10
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func NewBraintreeService(merchID, pubKey, privKey string) (Service, error) {
	// Create retry transport for retrying on io.EOF errors.
	transport := rehttp.NewTransport(nil,
		func(attempt rehttp.Attempt) bool {
			if attempt.Error == nil || attempt.Index == MaxTries {
				return false
			}

			return true
		},
		rehttp.ConstDelay(time.Millisecond*100),
	)

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   MaxTimeout,
	}

	api := braintree.NewWithHttpClient(braintree.Sandbox, merchID, pubKey, privKey, httpClient)

	return &braintreeService{
		account: account{
			merchID: merchID,
			pubKey:  pubKey,
			privKey: privKey,
		},
		api: api,
	}, nil
}

func (s *braintreeService) Type() ServiceType {
	return braintreeServiceType
}

func (s *braintreeService) BufferSize() int {
	return 1800
}

func (s *braintreeService) MaxThreads() int {
	return 15
}

func (s *braintreeService) NewBackup() error {
	return nil
}

func (s *braintreeService) Upload(data []byte) (*ServiceData, error) {
	cf := make(map[string]string)
	for i := 1; i <= NumCustomFields; i++ {
		start := (i - 1) * len(data) / NumCustomFields
		end := i * len(data) / NumCustomFields

		if start > len(data) {
			break
		}

		if end > len(data) {
			end = len(data)
		}

		cf[fmt.Sprintf("data%d", i)] = string(data[start:end])
	}

	c, err := s.api.Customer().Create(context.Background(), &braintree.CustomerRequest{
		Email:        randomEmail(),
		CustomFields: cf,
	})
	if err != nil {
		return nil, err
	}

	return &ServiceData{
		Data: []byte(c.Id),
	}, nil
}

func (s *braintreeService) NewDownload() error {
	return nil
}

func (s *braintreeService) Download(sd ServiceData) ([]byte, error) {
	c, err := s.api.Customer().Find(context.Background(), string(sd.Data))
	if err != nil {
		return nil, err
	}

	data := []byte{}
	for i := 1; i <= NumCustomFields; i++ {
		var cfData string
		var ok bool
		if cfData, ok = c.CustomFields[fmt.Sprintf("data%d", i)]; !ok {
			break
		}

		data = append(data, []byte(cfData)...)
	}

	return data, nil
}

func randomEmail() string {
	n := rand.Intn(20) + 5
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b) + "@email.com"
}
