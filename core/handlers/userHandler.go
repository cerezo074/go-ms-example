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
		return handler.send(nil, &response.ResponseError{StatusCode: http.StatusInternalServerError, Message: err.Error()}, context)
	}

	if len(users) == 0 {
		return handler.send([]entities.User{}, nil, context)
	}

	return handler.send(users, nil, context)
}

func (handler userHandler) getUser(context *fiber.Ctx) error {
	userEmail := context.Query("address")

	if userEmail == "" {
		return handler.send(nil, &response.ResponseError{StatusCode: http.StatusBadRequest, Message: "address is not present on url as a query param"}, context)
	}

	user, err := handler.store.User(userEmail)

	if err != nil {
		return handler.send(nil, &response.ResponseError{StatusCode: http.StatusNotFound, Message: err.Error()}, context)
	}

	return handler.send(user, nil, context)
}

func (handler userHandler) newUser(context *fiber.Ctx) error {
	// user := new(entities.User)

	// if err := context.BodyParser(user); err != nil {
	// 	context.SendStatus(503)
	// 	return nil
	// }

	// return context.JSON(fiber.Map{
	// 	"result": fmt.Sprintf("Welcome %s!", user.Nickname),
	// })

	return context.SendStatus(http.StatusMethodNotAllowed)
}

func (handler userHandler) updateUser(context *fiber.Ctx) error {
	// //USE Multipart form data and get the url from s3
	// userEmail := context.Query("address")
	// if userEmail == "" {
	// 	return handler.send(nil, &response.ResponseError{StatusCode: http.StatusBadRequest, Message: "address is not present on url as a query param"}, context)
	// }

	// updatedUser := new(entities.User)
	// if err := context.BodyParser(updatedUser); err != nil {
	// 	return handler.send(nil, &response.ResponseError{StatusCode: http.StatusBadRequest, Message: err.Error()}, context)
	// }

	// //USE Multipart form data and get the url from s3, when we get the url we set it on updatedUser
	// oldUserValue, err := handler.store.User(userEmail)
	// oldUser := &oldUserValue
	// oldUser.Update(*updatedUser)
	// invalidUser := &response.ResponseError{StatusCode: http.StatusNoContent, Message: "user doesn't exists"}

	// if err == nil {
	// 	err = handler.store.UpdateUser(oldUser)
	// 	if err != nil {
	// 		return handler.send(nil, invalidUser, context)
	// 	}

	// 	return handler.send("user deleted successfully", nil, context)
	// }

	// return handler.send(nil, invalidUser, context)
	return context.SendStatus(http.StatusMethodNotAllowed)
}

func (handler userHandler) deleteUser(context *fiber.Ctx) error {
	userEmail := context.Query("address")

	if userEmail == "" {
		return handler.send(nil, &response.ResponseError{StatusCode: http.StatusBadRequest, Message: "address is not present on url as a query param"}, context)
	}

	err := handler.store.DeleteUser(userEmail)

	if err != nil {
		return handler.send(nil, &response.ResponseError{StatusCode: http.StatusNoContent, Message: "user doesn't exists"}, context)
	}

	return handler.send("user deleted successfully", nil, context)
}

func (handler userHandler) send(value interface{}, err *response.ResponseError, context *fiber.Ctx) error {
	if err != nil {
		return response.MakeJSON(response.Fail, nil, err, context)
	}

	return response.MakeJSON(response.Success, &value, nil, context)
}
