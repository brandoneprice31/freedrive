package config

type (
	Config struct {
		Email    string
		Password string
		KeyPath  string
	}
)

func New(email, password, keypath string) Config {
	return Config{
		Email:    email,
		Password: password,
		KeyPath:  keypath,
	}
}
