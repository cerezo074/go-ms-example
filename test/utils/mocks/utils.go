package mocks

import (
	"user/app/utils/config"
	"user/core/dependencies/services"
	. "user/core/dependencies/services"
	"user/core/entities"
)

func NewUserMockedServices(userRepository entities.UserRepository, userValidator UserValidatorServices, userImage ProfileImageServices) services.App {
	fakeConfig := config.Credentials{}

	fakeRepo := RepositoryServices{
		UserRepository: userRepository,
	}

	fakeValidor := ValidatorServices{
		UserValidator: userValidator,
	}

	fakeImage := ImageServices{
		UserProfileImage: userImage,
	}

	return App{
		Credentials: fakeConfig,
		Repository:  fakeRepo,
		Validator:   fakeValidor,
		Image:       fakeImage,
	}
}
