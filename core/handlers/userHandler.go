package handlers

import (
	"bytes"
	"net/http"
	"user/app/utils/config"
	"user/app/utils/response"
	"user/core/entities"
	"user/core/middleware/amazons3"
	"user/core/middleware/validator"

	"github.com/gofiber/fiber/v2"
)

//NewUserHandler factory method for creating a method handler for users
func NewUserHandler(repository entities.UserRepository, config config.Credentials) MethodHandlers {
	return userHandler{store: repository, credentials: config}
}

const (
	imagePath = "/api/v1/users/image/"
)

type userHandler struct {
	store       entities.UserRepository
	credentials config.Credentials
}

type BasicUser struct {
	Name  string
	Email string
}

func (handler userHandler) RegisterMethods(app *fiber.App) {
	app.Get("/api/v1/users", handler.getUsers)
	app.Get("/api/v1/users/email", handler.getUser)
	app.Get(imagePath+":id", amazons3.NewDownloader(handler.credentials), handler.getImage)
	app.Post("/api/v1/users", validator.DuplicatedUser(handler.store), amazons3.NewUploader(handler.credentials), handler.newUser)
	app.Put("/api/v1/users", handler.updateUser)
	app.Delete("/api/v1/users/email", handler.deleteUser)
}

func (handler userHandler) getUsers(context *fiber.Ctx) error {
	users, err := handler.store.Users()

	if err != nil {
		return response.MakeErrorJSON(http.StatusInternalServerError, err.Error())
	}

	if len(users) == 0 {
		return response.MakeSuccessJSON([]entities.User{}, context)
	}

	return response.MakeSuccessJSON(users, context)
}

func (handler userHandler) getUser(context *fiber.Ctx) error {
	userEmail := context.Query("address")

	if userEmail == "" {
		return response.MakeErrorJSON(http.StatusBadRequest, "address is not present on url as a query param")
	}

	user, err := handler.store.User(userEmail)

	if err != nil {
		return response.MakeErrorJSON(http.StatusNotFound, err.Error())
	}

	return response.MakeSuccessJSON(user, context)
}

func (handler userHandler) getImage(context *fiber.Ctx) error {
	if s3DataFile, ok := context.Locals(amazons3.S3_DOWNLOADED_IMAGE_FILE).(*amazons3.AWSBufferedFile); ok {
		return context.Status(http.StatusOK).SendStream(bytes.NewReader(s3DataFile.Data), int(s3DataFile.Size))
	}

	return response.MakeErrorJSON(http.StatusInternalServerError, "Invalid type file")
}

func (handler userHandler) newUser(context *fiber.Ctx) error {
	user := new(entities.User)

	if err := context.BodyParser(user); err != nil {
		return response.MakeErrorJSON(http.StatusNotFound, err.Error())
	}

	if imageURI, ok := context.Locals(amazons3.S3_UPLOADED_IMAGE_ID).(string); ok {
		user.ImageID = imagePath + imageURI
	} else {
		user.ImageID = imagePath + amazons3.DEFAULT_IMAGE
	}

	if err := handler.store.CreateUser(user); err != nil {
		return response.MakeErrorJSON(http.StatusBadRequest, err.Error())
	}

	return response.MakeSuccessJSON("user was created successfully", context)
}

func (handler userHandler) updateUser(context *fiber.Ctx) error {
	updatedUser := new(entities.User)
	if err := context.BodyParser(updatedUser); err != nil {
		return response.MakeErrorJSON(http.StatusBadRequest, err.Error())
	}

	if err := updatedUser.IsValid(); err != nil {
		return response.MakeErrorJSON(http.StatusBadRequest, err.Error())
	}

	if err := handler.store.UpdateUser(updatedUser); err != nil {
		return response.MakeErrorJSON(http.StatusInternalServerError, err.Error())
	}

	return response.MakeSuccessJSON("user updated successfully", context)
}

func (handler userHandler) deleteUser(context *fiber.Ctx) error {
	userEmail := context.Query("address")

	if userEmail == "" {
		return response.MakeErrorJSON(http.StatusBadRequest, "address is not present on url as a query param")
	}

	if err := handler.store.DeleteUser(userEmail); err != nil {
		return response.MakeErrorJSON(http.StatusInternalServerError, err.Error())
	}

	return response.MakeSuccessJSON("user deleted successfully", context)
}
