package dependency

import (
	"user/app/utils/config"
	"user/core/dependencies/services"
)

func NewServiceLocator(credentials *config.Credentials) (*services.App, error) {
	appCredentials := credentials
	if appCredentials == nil {
		defaultCredentials, err := NewCredentials(DefaultCredentialsPath)
		if err != nil {
			return nil, err
		}

		appCredentials = defaultCredentials
	}

	appRepository, err := NewRepository(*appCredentials)
	if err != nil {
		return nil, err
	}

	appValidator := NewValidator(*appRepository)
	appImages := NewImageProvider(*appRepository, appValidator, *appCredentials, AWSS3)

	return &services.App{
		Credentials: *appCredentials,
		Repository:  *appRepository,
		Validator:   appValidator,
		Image:       appImages,
	}, nil
}
