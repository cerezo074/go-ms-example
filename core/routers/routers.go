package routers

import (
	"user/app/utils/config"
	"user/core/entities"

	"github.com/gofiber/fiber/v2"
)

//RouteHandler declares an abstract interface for using with specific routers
type RouteHandler interface {
	Register(app *fiber.App, repository entities.Repository, credentials config.Credentials)
}
