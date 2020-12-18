package handlers

import (
	"fmt"
	"user/core/entities"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
}

//output
type BasicUser struct {
	Name  string
	Email string
}

//BuildUserHandler factory method for creating a method handler to users
func BuildUserHandler() MethodHandlers {
	return userHandler{}
}

func (handler userHandler) RegisterMethods(app *fiber.App) {
	app.Get("/api/v1/users", handler.getUsers)
	app.Get("/api/v1/users/:id", handler.getUser)
	app.Post("/api/v1/users", handler.newUser)
	app.Put("/api/v1/users/:id", handler.updateUser)
	app.Delete("/api/v1/users/:id", handler.deleteUser)
}

func (handler userHandler) getUsers(c *fiber.Ctx) error {
	users := []BasicUser{
		{Name: "Usuario 1", Email: "Email1@g.co"},
		{Name: "Usuario 2", Email: "Email2@g.co"},
	}

	return c.JSON(fiber.Map{
		"result": users,
	})
}

func (handler userHandler) getUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	if userID == "" {
		c.SendStatus(503)
		return nil
	}

	user := BasicUser{
		Name:  fmt.Sprintf("Usuario %s!", userID),
		Email: fmt.Sprintf("email%s@f.co", userID),
	}

	return c.JSON(fiber.Map{
		"result": user,
	})
}

func (handler userHandler) newUser(c *fiber.Ctx) error {
	user := new(entities.User)

	if err := c.BodyParser(user); err != nil {
		c.SendStatus(503)
		return nil
	}

	return c.JSON(fiber.Map{
		"result": fmt.Sprintf("Welcome %s!", user.Name),
	})
}

func (handler userHandler) updateUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	if userID == "" {
		c.SendStatus(503)
		return nil
	}

	return c.JSON(fiber.Map{
		"result": fmt.Sprintf("User %s updated!", userID),
	})
}

func (handler userHandler) deleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	if userID == "" {
		c.SendStatus(503)
		return nil
	}

	return c.JSON(fiber.Map{
		"result": fmt.Sprintf("User %s deleted!", userID),
	})
}
