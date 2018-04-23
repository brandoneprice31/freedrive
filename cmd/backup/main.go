package main

import (
	"github.com/brandoneprice31/freedrive/setup"
)

func main() {
	m, err := setup.Manager()
	if err != nil {
		panic(err)
	}

	m.Backup()
}
