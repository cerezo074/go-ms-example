package services

import "github.com/gofiber/fiber/v2"

type ValidatorServices struct {
	UserValidator UserValidatorServices
}

type UserValidatorServices interface {
	DuplicatedUser() fiber.Handler
	IsValidEmailFormat(email string) bool
}
