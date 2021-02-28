package dependency

import (
	"user/app/db/postgres"
	"user/app/utils/config"
	"user/core/entities"
	"user/core/services"
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

// type FakeUserRepository struct {
// }

// func (repo FakeUserRepository) User(email string) (entities.User, error) {

// }

// func (repo FakeUserRepository) Users() ([]entities.User, error) {

// }

// func (repo FakeUserRepository) CreateUser(user *entities.User) error {

// }

// func (repo FakeUserRepository) UpdateUser(user *entities.User) error {

// }

// func (repo FakeUserRepository) DeleteUser(email string) error {

// }

// func (repo FakeUserRepository) ExistUser(email string) bool {

// }

// func FakeRepository() entities.UserRepository {

// }
