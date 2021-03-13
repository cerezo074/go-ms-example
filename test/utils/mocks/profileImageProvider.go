package mocks

import (
	"github.com/gofiber/fiber/v2"
)

type FakeProfileImage struct {
	Delete func(*fiber.Ctx) error
}

func (object FakeProfileImage) NewUploader() fiber.Handler {
	return func(context *fiber.Ctx) error {
		return context.Next()
	}
}

func (object FakeProfileImage) NewDownloader() fiber.Handler {
	return func(context *fiber.Ctx) error {
		return context.Next()
	}
}

func (object FakeProfileImage) DeleteImage() fiber.Handler {
	return func(context *fiber.Ctx) error {
		if object.Delete == nil {
			return context.Next()
		}

		return object.Delete(context)
	}
}

func (object FakeProfileImage) UpdateImage() fiber.Handler {
	return func(context *fiber.Ctx) error {
		return context.Next()
	}
}
