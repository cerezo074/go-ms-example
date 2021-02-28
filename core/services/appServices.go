package services

import "user/app/utils/config"

type App struct {
	Credentials config.Credentials
	Repository  RepositoryServices
	Validator   ValidatorServices
	Image       ImageServices
}
