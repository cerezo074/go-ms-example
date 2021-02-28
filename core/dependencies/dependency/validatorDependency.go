package dependency

import (
	"user/core/dependencies/services"
	"user/core/middleware/validator"
)

func NewValidator(repository services.RepositoryServices) services.ValidatorServices {
	userValidator := validator.UserValidatorProvider{
		UserStore: repository.UserRepository,
	}

	return services.ValidatorServices{
		UserValidator: userValidator,
	}
}
