package dependency

import (
	"fmt"
	"user/app/utils/config"
	"user/core/dependencies/services"
	"user/core/middleware/image"
)

type LoaderType int

const (
	AWSS3 LoaderType = iota
)

type ImageStorageBuilder struct {
	LoaderType LoaderType
}

func (object ImageStorageBuilder) Load(credentials config.Credentials) (services.ImageStorageSession, error) {
	switch object.LoaderType {
	case AWSS3:
		return image.NewS3StorageSession(credentials)
	default:
		return nil, fmt.Errorf("Invalid image loader type for value %v, ", object.LoaderType)
	}
}

func NewImageProvider(
	repository services.RepositoryServices,
	validator services.ValidatorServices,
	credentials config.Credentials,
	imageLoader LoaderType) services.ImageServices {

	imageProvider := image.ProfileImageProvider{
		Credentials:   credentials,
		UserStore:     repository.UserRepository,
		UserValidator: validator.UserValidator,
		Loader:        ImageStorageBuilder{LoaderType: imageLoader},
	}

	return services.ImageServices{
		UserProfileImage: imageProvider,
	}
}
