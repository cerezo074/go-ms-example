package dependency

import (
	"user/app/utils/config"
)

var (
	DefaultCredentialsPath = "../app/"
)

func NewCredentials(path string) (*config.Credentials, error) {
	config, err := config.LoadCredentials(path)
	return &config, err
}

func FakeCredentials() config.Credentials {
	return config.Credentials{}
}
