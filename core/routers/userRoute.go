package routers

import (
	"user/core/handlers"
	"user/core/services"

	"github.com/gofiber/fiber/v2"
)

//NewUserRouter factory method for creating a route handler to users
func NewUserRouter() RouteHandler {
	return userRoutes{}
}

type userRoutes struct {
}

//Register insert user handlers into app
func (router userRoutes) Register(app *fiber.App, appDependencies services.App) {
	crudHandler := handlers.NewUserHandler(appDependencies)
	crudHandler.RegisterMethods(app)
}
