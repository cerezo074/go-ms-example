package routers

import (
	"user/core/entities"
	"user/core/handlers"

	"github.com/gofiber/fiber/v2"
)

//NewUserRouter factory method for creating a route handler to users
func NewUserRouter() RouteHandler {
	return userRoutes{}
}

type userRoutes struct {
}

//Register insert user handlers into app
func (router userRoutes) Register(app *fiber.App, repository entities.Repository) {
	crudHandler := handlers.NewUserHandler(repository)
	crudHandler.RegisterMethods(app)
}