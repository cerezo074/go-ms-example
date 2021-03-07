package mocks

import (
	"user/core/entities"

	"github.com/gofiber/fiber/v2"
)

type FakeValidator struct {
}

func (object FakeValidator) DuplicatedUser() fiber.Handler {
	return func(context *fiber.Ctx) error {
		return context.Next()
	}
}

func (object FakeValidator) IsValidEmailFormat(email string) bool {
	return true
}

func (object FakeValidator) IsValid(user entities.User) error {
	return nil
}
