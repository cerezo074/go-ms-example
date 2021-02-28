package dependency

import (
	"user/app/utils/config"
	"user/core/dependencies/services"
	"user/core/middleware/amazons3"
)

func NewImageProvider(
	repository services.RepositoryServices,
	validator services.ValidatorServices,
	credentials config.Credentials) services.ImageServices {
	imageProvider := amazons3.S3ProfileImageProvider{
		Credentials:   credentials,
		UserStore:     repository.UserRepository,
		UserValidator: validator.UserValidator,
	}

	return services.ImageServices{
		UserProfileImage: imageProvider,
	}
}
