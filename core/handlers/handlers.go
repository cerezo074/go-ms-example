package handlers

import "github.com/gofiber/fiber/v2"

type MethodHandlers interface {
	RegisterMethods(app *fiber.App)
}
