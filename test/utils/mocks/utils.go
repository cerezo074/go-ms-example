package mocks

import (
	"user/app/utils/config"
	"user/core/dependencies/services"
	. "user/core/dependencies/services"
	"user/core/entities"
	"user/core/middleware/image"
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

func NewImageProvider(userRepository entities.UserRepository,
	userValidator UserValidatorServices,
	imageLoder ImageStorageLoader) image.ProfileImageProvider {

	return image.ProfileImageProvider{
		Credentials:   fakeCredentias,
		UserStore:     userRepository,
		UserValidator: userValidator,
		Loader:        imageLoder,
	}
}
