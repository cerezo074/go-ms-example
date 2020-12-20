package handlers

import (
	"fmt"
	"user/app/utils/response"
	"user/core/entities"

	"github.com/gofiber/fiber/v2"
)

//NewUserHandler factory method for creating a method handler to users
func NewUserHandler(repository entities.UserRepository) MethodHandlers {
	return userHandler{store: repository}
}

type userHandler struct {
	store entities.UserRepository
}

type BasicUser struct {
	Name  string
	Email string
}

func (handler userHandler) RegisterMethods(app *fiber.App) {
	app.Get("/api/v1/users", handler.getUsers)
	app.Get("/api/v1/users/:id", handler.getUser)
	app.Post("/api/v1/users", handler.newUser)
	app.Put("/api/v1/users/:id", handler.updateUser)
	app.Delete("/api/v1/users/:id", handler.deleteUser)
}

func (handler userHandler) getUsers(context *fiber.Ctx) error {
	users, err := handler.store.Users()

	if err != nil {
		return response.MakeJSON(response.Fail, nil, err, context)
	}

	return response.MakeJSON(response.Success, users, nil, context)
}

func (handler userHandler) getUser(context *fiber.Ctx) error {
	userID := context.Params("id")

	if userID == "" {
		context.SendStatus(503)
		return nil
	}

	user := BasicUser{
		Name:  fmt.Sprintf("Usuario %s!", userID),
		Email: fmt.Sprintf("email%s@f.co", userID),
	}

	return context.JSON(fiber.Map{
		"result": user,
	})
}

func (handler userHandler) newUser(context *fiber.Ctx) error {
	user := new(entities.User)

	if err := context.BodyParser(user); err != nil {
		context.SendStatus(503)
		return nil
	}

	return context.JSON(fiber.Map{
		"result": fmt.Sprintf("Welcome %s!", user.Nickname),
	})
}

func (handler userHandler) updateUser(context *fiber.Ctx) error {
	userID := context.Params("id")

	if userID == "" {
		context.SendStatus(503)
		return nil
	}

	return context.JSON(fiber.Map{
		"result": fmt.Sprintf("User %s updated!", userID),
	})
}

func (handler userHandler) deleteUser(context *fiber.Ctx) error {
	userID := context.Params("id")

	if userID == "" {
		context.SendStatus(503)
		return nil
	}

	return context.JSON(fiber.Map{
		"result": fmt.Sprintf("User %s deleted!", userID),
	})
}
