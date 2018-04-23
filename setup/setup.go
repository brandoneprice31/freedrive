package setup

import (
	"fmt"

	"github.com/brandoneprice31/freedrive/config"
	"github.com/brandoneprice31/freedrive/manager"
	"github.com/brandoneprice31/freedrive/service"
)

func Manager(fd string) (*manager.Manager, error) {
	c := config.New("", "", fmt.Sprintf("%s/key", fd))

	bts, err := service.NewBraintreeService("q9t77sxx3jngp9qb", "s99bbz4mg8qqf4b4", "39f09de6aca920b82ffa982664b7fbaf")
	if err != nil {
		return nil, err
	}

	dbs, err := service.NewDropboxService("pYTxO8EdjVEAAAAAAAAF_JqkFisAR6HLpjNlBSy1crQ_xtw1aTMiHx5aS0VV4UgW")
	if err != nil {
		return nil, err
	}

	tws, err := service.NewTwitterService("988142103086751744-JvWk6EHpy1O0WDYkNvUhKqKeztdlrSm", "fqQwBgtyJBoyyCqH6GpogP4OBJ4O7RIhBc5gdzhcTFxPR", "ex1gHG1AFYCf8powfJAofgU4j", "qtyge4EfEtA9vBzGguLMX30HPVrZhUu0yni1i56E8N7eYQGyim")
	if err != nil {
		return nil, err
	}

	return manager.NewManager(c, bts, dbs, tws), nil
}
