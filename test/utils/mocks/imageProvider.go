package mocks

import "github.com/gofiber/fiber/v2"

type FakeImage struct {
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
		return context.Next()
	}
}
func (object FakeImage) UpdateImage() fiber.Handler {
	return func(context *fiber.Ctx) error {
		return context.Next()
	}
}
