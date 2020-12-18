package routers

import (
	"user/core/handlers"

	"github.com/gofiber/fiber/v2"
)

type userRoutes struct {
}

//BuildUserRouter factory method for creating a route handler to users
func BuildUserRouter() RouteHandler {
	return userRoutes{}
}

//Register insert user handlers into app
func (router userRoutes) Register(app *fiber.App) {
	crudHandler := handlers.BuildUserHandler()
	crudHandler.RegisterMethods(app)
}
