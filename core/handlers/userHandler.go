package handlers

import (
	"net/http"
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
	app.Get("/api/v1/users/email", handler.getUser)
	app.Post("/api/v1/users", handler.newUser)
	app.Put("/api/v1/users", handler.updateUser)
	app.Delete("/api/v1/users/email", handler.deleteUser)
}

func (handler userHandler) getUsers(context *fiber.Ctx) error {
	users, err := handler.store.Users()

	if err != nil {
		return handler.sendError(http.StatusInternalServerError, err.Error())
	}

	if len(users) == 0 {
		return response.MakeSuccessJSON([]entities.User{}, context)
	}

	return response.MakeSuccessJSON(users, context)
}

func (handler userHandler) getUser(context *fiber.Ctx) error {
	userEmail := context.Query("address")

	if userEmail == "" {
		return handler.sendError(http.StatusBadRequest, "address is not present on url as a query param")
	}

	user, err := handler.store.User(userEmail)

	if err != nil {
		return handler.sendError(http.StatusNotFound, err.Error())
	}

	return response.MakeSuccessJSON(user, context)
}

func (handler userHandler) newUser(context *fiber.Ctx) error {
	user := new(entities.User)

	if err := context.BodyParser(user); err != nil {
		return handler.sendError(http.StatusNotFound, err.Error())
	}

	if err := handler.store.CreateUser(user); err != nil {
		return handler.sendError(http.StatusBadRequest, err.Error())
	}

	return response.MakeSuccessJSON("user was created successfully", context)
}

func (handler userHandler) updateUser(context *fiber.Ctx) error {
	updatedUser := new(entities.User)
	if err := context.BodyParser(updatedUser); err != nil {
		return handler.sendError(http.StatusBadRequest, err.Error())
	}

	if err := updatedUser.IsValid(); err != nil {
		return handler.sendError(http.StatusBadRequest, err.Error())
	}

	if err := handler.store.UpdateUser(updatedUser); err != nil {
		return handler.sendError(http.StatusInternalServerError, err.Error())
	}

	return response.MakeSuccessJSON("user updated successfully", context)
}

func (handler userHandler) deleteUser(context *fiber.Ctx) error {
	userEmail := context.Query("address")

	if userEmail == "" {
		return handler.sendError(http.StatusBadRequest, "address is not present on url as a query param")
	}

	if err := handler.store.DeleteUser(userEmail); err != nil {
		return handler.sendError(http.StatusInternalServerError, err.Error())
	}

	return response.MakeSuccessJSON("user deleted successfully", context)
}

func (handler userHandler) sendError(httpStatusCode int, description string) error {
	return response.ResponseError{
		StatusCode: httpStatusCode,
		Message:    description,
	}
}
