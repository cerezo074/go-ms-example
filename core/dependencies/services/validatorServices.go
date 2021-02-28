package services

import (
	"user/core/entities"

	"github.com/gofiber/fiber/v2"
)

type ValidatorServices struct {
	UserValidator UserValidatorServices
}

type UserValidatorServices interface {
	DuplicatedUser() fiber.Handler
	IsValidEmailFormat(email string) bool
	IsValid(user entities.User) error
}
