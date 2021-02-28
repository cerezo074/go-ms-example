package dependency

import (
	"user/core/middleware/validator"
	"user/core/services"
)

func NewValidator(repository services.RepositoryServices) services.ValidatorServices {
	userValidator := validator.UserValidatorProvider{
		UserStore: repository.UserRepository,
	}

	return services.ValidatorServices{
		UserValidator: userValidator,
	}
}
