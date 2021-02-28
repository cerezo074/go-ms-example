package handlers

import (
	"bytes"
	"net/http"
	"user/app/utils/config"
	"user/app/utils/response"
	"user/core/entities"
	"user/core/middleware/amazons3"
	"user/core/services"

	"github.com/gofiber/fiber/v2"
)

//NewUserHandler factory method for creating a method handler for users
func NewUserHandler(depedencies services.App) MethodHandlers {
	return userHandler{appDependencies: depedencies}
}

const (
	imagePath = "/api/v1/users/image/"
)

type userHandler struct {
	appDependencies services.App
}

func (object userHandler) userImage() services.S3ProfileImageServices {
	return object.appDependencies.Image.UserProfileImage
}

func (object userHandler) userValidator() services.UserValidatorServices {
	return object.appDependencies.Validator.UserValidator
}

func (object userHandler) appCredentials() config.Credentials {
	return object.appDependencies.Credentials
}

func (object userHandler) userRepository() entities.UserRepository {
	return object.appDependencies.Repository.UserRepository
}

type BasicUser struct {
	Name  string
	Email string
}

func (object userHandler) RegisterMethods(app *fiber.App) {
	app.Get("/api/v1/users", object.getUsers)
	app.Get("/api/v1/users/email", object.getUser)
	app.Get(imagePath+":id", object.userImage().NewDownloader(), object.getImage)
	app.Post("/api/v1/users", object.userValidator().DuplicatedUser(), object.userImage().NewUploader(), object.newUser)
	app.Put("/api/v1/users", object.userImage().UpdateImage(), object.updateUser)
	app.Delete("/api/v1/users/email", object.userImage().DeleteImage(), object.deleteUser)
}

func (object userHandler) getUsers(context *fiber.Ctx) error {
	users, err := object.userRepository().Users()

	if err != nil {
		return response.MakeErrorJSON(http.StatusInternalServerError, err.Error())
	}

	if len(users) == 0 {
		return response.MakeSuccessJSON([]entities.User{}, context)
	}

	return response.MakeSuccessJSON(users, context)
}

func (object userHandler) getUser(context *fiber.Ctx) error {
	userEmail := context.Query("address")

	if userEmail == "" {
		return response.MakeErrorJSON(http.StatusBadRequest, "address is not present on url as a query param")
	}

	user, err := object.userRepository().User(userEmail)

	if err != nil {
		return response.MakeErrorJSON(http.StatusNotFound, err.Error())
	}

	return response.MakeSuccessJSON(user, context)
}

func (object userHandler) getImage(context *fiber.Ctx) error {
	if s3DataFile, ok := context.Locals(amazons3.S3_DOWNLOADED_IMAGE_FILE).(*amazons3.AWSBufferedFile); ok {
		return context.Status(http.StatusOK).SendStream(bytes.NewReader(s3DataFile.Data), int(s3DataFile.Size))
	}

	return response.MakeErrorJSON(http.StatusInternalServerError, "Invalid type file")
}

func (object userHandler) newUser(context *fiber.Ctx) error {
	user := new(entities.User)

	if err := context.BodyParser(user); err != nil {
		return response.MakeErrorJSON(http.StatusNotFound, err.Error())
	}

	if imageURI, ok := context.Locals(amazons3.S3_UPLOADED_IMAGE_ID).(string); ok {
		user.ImageID = imagePath + imageURI
	} else {
		user.ImageID = imagePath + amazons3.DEFAULT_IMAGE
	}

	if err := object.userRepository().CreateUser(user); err != nil {
		return response.MakeErrorJSON(http.StatusBadRequest, err.Error())
	}

	return response.MakeSuccessJSON("user was created successfully", context)
}

func (object userHandler) updateUser(context *fiber.Ctx) error {
	updatedUser := new(entities.User)
	if err := context.BodyParser(updatedUser); err != nil {
		return response.MakeErrorJSON(http.StatusBadRequest, err.Error())
	}

	if imageID, ok := context.Locals(amazons3.S3_UPLOADED_IMAGE_ID).(string); ok {
		updatedUser.ImageID = imagePath + imageID
	} else {
		updatedUser.ImageID = imagePath + amazons3.DEFAULT_IMAGE
	}

	if err := object.userValidator().IsValid(*updatedUser); err != nil {
		return response.MakeErrorJSON(http.StatusBadRequest, err.Error())
	}

	if err := object.userRepository().UpdateUser(updatedUser); err != nil {
		return response.MakeErrorJSON(http.StatusInternalServerError, err.Error())
	}

	return response.MakeSuccessJSON("user updated successfully", context)
}

func (object userHandler) deleteUser(context *fiber.Ctx) error {
	var user entities.User

	if assertion, ok := context.Locals(amazons3.S3_USER_ENTITY).(entities.User); ok {
		user = assertion
	}

	if err := object.userRepository().DeleteUser(user.Email); err != nil {
		return response.MakeErrorJSON(http.StatusInternalServerError, "Invalid user")
	}

	return response.MakeSuccessJSON("user deleted successfully", context)
}
