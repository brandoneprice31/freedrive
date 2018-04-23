package setup

import (
	"github.com/brandoneprice31/freedrive/config"
	"github.com/brandoneprice31/freedrive/manager"
	"github.com/brandoneprice31/freedrive/service"
)

func Manager() (*manager.Manager, error) {
	c := config.New("", "")

	bts, err := service.NewBraintreeService(c.Braintree.MerchantID, c.Braintree.PublicKey, c.Braintree.PrivateKey)
	if err != nil {
		return nil, err
	}

	dbs, err := service.NewDropboxService(c.Dropbox.AccessToken)
	if err != nil {
		return nil, err
	}

	tws, err := service.NewTwitterService(c.Twitter.AccessToken, c.Twitter.AccessSecret, c.Twitter.ConsumerKey, c.Twitter.ConsumerSecret)
	if err != nil {
		return nil, err
	}

	return manager.NewManager(c, bts, dbs, tws), nil
}
