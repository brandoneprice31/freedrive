package main

import (
	"fmt"
	"os"

	"github.com/brandoneprice31/freedrive/config"
	"github.com/brandoneprice31/freedrive/manager"
	"github.com/brandoneprice31/freedrive/service"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println(`3 args are required:
 - path/to/freedrive/dir
 - path/to/download/folder
 		`)
		os.Exit(2)
	}

	fd := os.Args[1]
	c := config.New("", "", fmt.Sprintf("%s/key", fd))

	_, err := service.NewLocalFileSystemService(fmt.Sprintf("%s/lfs", fd))
	if err != nil {
		panic(err)
	}

	bts, err := service.NewBraintreeService("q9t77sxx3jngp9qb", "s99bbz4mg8qqf4b4", "39f09de6aca920b82ffa982664b7fbaf")
	if err != nil {
		panic(err)
	}

	m := manager.NewManager(c, bts)

	downloadPath := os.Args[2]
	m.Download(downloadPath)
}
