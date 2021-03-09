package mocks

import (
	"github.com/gofiber/fiber/v2"
)

type FakeImage struct {
	Delete func(*fiber.Ctx) error
}

func (object FakeImage) NewUploader() fiber.Handler {
	return func(context *fiber.Ctx) error {
		return context.Next()
	}
}

func (object FakeImage) NewDownloader() fiber.Handler {
	return func(context *fiber.Ctx) error {
		return context.Next()
	}
}

func (object FakeImage) DeleteImage() fiber.Handler {
	return func(context *fiber.Ctx) error {
		if object.Delete == nil {
			return context.Next()
		}

		return object.Delete(context)
	}
}

func (object FakeImage) UpdateImage() fiber.Handler {
	return func(context *fiber.Ctx) error {
		return context.Next()
	}
}
