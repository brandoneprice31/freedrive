package config

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type (
	Config struct {
		Email    string
		Password string

		Paths     Paths     `yaml:"paths"`
		Braintree Braintree `yaml:"braintree"`
		Dropbox   Dropbox   `yaml:"dropbox"`
		Twitter   Twitter   `yaml:"twitter"`
	}

	Paths struct {
		Freedrive string `yaml:"freedrive"`
		Backup    string `yaml:"backup"`
	}

	Braintree struct {
		MerchantID string `yaml:"merchant_id"`
		PublicKey  string `yaml:"public_key"`
		PrivateKey string `yaml:"private_key"`
	}

	Dropbox struct {
		AccessToken string `yaml:"access_token"`
	}

	Twitter struct {
		AccessToken    string `yaml:"access_token"`
		AccessSecret   string `yaml:"access_secret"`
		ConsumerKey    string `yaml:"consumer_key"`
		ConsumerSecret string `yaml:"consumer_secret"`
	}
)

func New(email, password string) Config {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var c Config
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	c.Paths.Freedrive = fmt.Sprintf("%s/key", c.Paths.Freedrive)
	c.Email = email
	c.Password = password

	return c
}
