package dependency

import (
	"user/app/utils/config"
	"user/core/dependencies/services"
	"user/core/middleware/image"
)

func NewImageProvider(
	repository services.RepositoryServices,
	validator services.ValidatorServices,
	credentials config.Credentials,
	imageLoader image.LoaderType) services.ImageServices {

	imageProvider := image.ProfileImageProvider{
		Credentials:   credentials,
		UserStore:     repository.UserRepository,
		UserValidator: validator.UserValidator,
		Loader:        image.ImageStorageBuilder{LoaderType: imageLoader},
	}

	return services.ImageServices{
		UserProfileImage: imageProvider,
	}
}
