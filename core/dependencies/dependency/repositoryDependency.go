package dependency

import (
	"user/app/db/postgres"
	"user/app/utils/config"
	"user/core/dependencies/services"
	"user/core/entities"
)

func NewRepository(credentials config.Credentials) (*services.RepositoryServices, error) {
	userRepo, err := loadUserRepository(credentials)
	if err != nil {
		return nil, err
	}

	return &services.RepositoryServices{
		UserRepository: userRepo,
	}, nil
}

func loadUserRepository(credentials config.Credentials) (entities.UserRepository, error) {
	store, err := postgres.NewStore(credentials.DBSource, credentials.DBDriver)
	return store, err
}
