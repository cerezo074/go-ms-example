package routers

import (
	"fmt"
	"net/http"
	"testing"
	"user/app/utils/response"
	. "user/core/entities"
	"user/core/middleware/amazons3"
	"user/core/routers"
	utils "user/test/utils/http"
	. "user/test/utils/mocks"
	"user/test/utils/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	user1 = User{
		ID:          uuid.New(),
		Email:       "user1@gmail.com",
		Nickname:    "CR7",
		Password:    "123456",
		ImageID:     "profile1.png",
		CountryCode: "USA",
		Birthday:    "07/01/2000",
	}
	user2 = User{
		ID:          uuid.New(),
		Email:       "user2@gmail.com",
		Nickname:    "Messi",
		Password:    "654321",
		ImageID:     "profile2.png",
		CountryCode: "USA",
		Birthday:    "08/01/2000",
	}
	users = []User{
		user1,
		user2,
	}
	allUsersRepo = FakeRepo{
		AllUsers: func() ([]User, error) {
			return users, nil
		},
	}
	findUserRepo = FakeRepo{
		UserByEmail: func(email string) (User, error) {
			if email == user1.Email {
				return user1, nil
			}
			return User{}, fmt.Errorf("User doesn't exist with %v email", email)
		},
	}
	deleteUserRepo = FakeRepo{
		Delete: func(email string) error {
			if email == user1.Email {
				return nil
			}
			return fmt.Errorf("User doesn't exist with %v email", email)
		},
	}
	deleteUserImage = FakeImage{
		Delete: func(context *fiber.Ctx) error {
			context.Locals(amazons3.S3_USER_ENTITY, user1)
			return context.Next()
		},
	}
	deleteInvalidUserImage = FakeImage{
		Delete: func(context *fiber.Ctx) error {
			context.Locals(amazons3.S3_USER_ENTITY, User{})
			return context.Next()
		},
	}
)

func buildServer(unmarshaller utils.ResponseUnmarshaller) *utils.FakeServer {
	return &utils.FakeServer{
		FiberApp: fiber.New(fiber.Config{
			ErrorHandler: response.HandleJSONError,
		}),
		Unmarshaller: unmarshaller,
	}
}

func Test_ShouldGetUsers_WhenExistInRepository(t *testing.T) {
	t.Parallel()
	server := buildServer(utils.NewAllUsersUnmarshaller)
	fakeRepo := allUsersRepo
	appServices := NewUserMockedServices(fakeRepo, FakeValidator{}, FakeImage{})

	routers.NewUserRouter().Register(server.FiberApp, appServices)
	response, object, _ := server.Execute("GET", "/api/v1/users", nil)
	userSlice, ok := object.([]models.User)

	assert.Truef(t, ok, "Invalid type marshalled from response")
	assert.EqualValuesf(t, http.StatusOK, response.StatusCode, "Invalid http status code")
	assert.Lenf(t, userSlice, len(users), "Invalid length for users, must be 2")
	assert.EqualValuesf(t, userSlice[0].Email, "user1@gmail.com", "Invalid email for user1")
	assert.EqualValuesf(t, userSlice[1].Email, "user2@gmail.com", "Invalid email for user2")
}

func Test_ShouldntGetUsers_WhenRepositoryIsEmpty(t *testing.T) {
	t.Parallel()
	server := buildServer(utils.NewAllUsersUnmarshaller)
	appServices := NewUserMockedServices(FakeRepo{}, FakeValidator{}, FakeImage{})

	routers.NewUserRouter().Register(server.FiberApp, appServices)
	response, object, _ := server.Execute("GET", "/api/v1/users", nil)
	userSlice, ok := object.([]models.User)

	assert.Truef(t, ok, fmt.Sprintf("Invalid type marshalled %t from response body, expected []User type", object))
	assert.EqualValuesf(t, http.StatusOK, response.StatusCode, "Invalid http status code")
	assert.Lenf(t, userSlice, 0, "Invalid length for users, must be 0")
}

func Test_ShouldGetUserByEmail_WhenExistInRepository(t *testing.T) {
	t.Parallel()
	server := buildServer(utils.NewFindUserUnmarshaller)
	fakeRepo := findUserRepo
	appServices := NewUserMockedServices(fakeRepo, FakeValidator{}, FakeImage{})

	routers.NewUserRouter().Register(server.FiberApp, appServices)
	response, object, _ := server.Execute("GET", "/api/v1/users/email?address="+user1.Email, nil)
	user, ok := object.(models.User)

	assert.Truef(t, ok, fmt.Sprintf("Invalid type marshalled %t from response body, expected User type", object))
	assert.EqualValuesf(t, http.StatusOK, response.StatusCode, "Invalid http status code")
	assert.EqualValuesf(t, user.Email, "user1@gmail.com", "Invalid email for user1")
	assert.EqualValuesf(t, user.Nickname, "CR7", "Invalid nickname for user1")
}

func Test_ShouldntGetUserByEmail_WhenUserDoesntExist(t *testing.T) {
	t.Parallel()
	server := buildServer(utils.NewFailUnmarshaller)
	fakeRepo := findUserRepo
	appServices := NewUserMockedServices(fakeRepo, FakeValidator{}, FakeImage{})
	invalidEmail := "invaliduser@test.com"

	routers.NewUserRouter().Register(server.FiberApp, appServices)
	response, object, _ := server.Execute("GET", "/api/v1/users/email?address="+invalidEmail, nil)
	failResponse, ok := object.(models.FailResponse)

	assert.Truef(t, ok, fmt.Sprintf("Invalid type marshalled %t from response body, expected FailResponse type", object))
	assert.EqualValuesf(t, http.StatusNotFound, response.StatusCode, "Invalid http status code")
	assert.EqualValuesf(t, failResponse.Error, fmt.Sprintf("User doesn't exist with %v email", invalidEmail), "Invalid error message")
}

func Test_DeleteUserByEmail_WhenExistInRepository(t *testing.T) {
	t.Parallel()
	server := buildServer(utils.NewSuccessUnmarshaller)
	fakeRepo := deleteUserRepo
	fakeImage := deleteUserImage
	appServices := NewUserMockedServices(fakeRepo, FakeValidator{}, fakeImage)

	routers.NewUserRouter().Register(server.FiberApp, appServices)
	response, object, _ := server.Execute("DELETE", "/api/v1/users/email?address="+user1.Email, nil)
	successResponse, ok := object.(models.SuccessResponse)

	assert.Truef(t, ok, fmt.Sprintf("Invalid type marshalled %t from response body, expected SuccessResponse type", object))
	assert.EqualValuesf(t, http.StatusOK, response.StatusCode, "Invalid http status code")
	assert.EqualValuesf(t, successResponse.Data, "user deleted successfully", "Invalid error message")
}

func Test_DeleteUserByEmail_WhenUserDoesntExist(t *testing.T) {
	t.Parallel()
	server := buildServer(utils.NewFailUnmarshaller)
	fakeRepo := deleteUserRepo
	fakeImage := deleteInvalidUserImage
	appServices := NewUserMockedServices(fakeRepo, FakeValidator{}, fakeImage)
	invalidEmail := "invaliduser@test.com"

	routers.NewUserRouter().Register(server.FiberApp, appServices)
	response, object, _ := server.Execute("DELETE", "/api/v1/users/email?address="+invalidEmail, nil)
	failResponse, ok := object.(models.FailResponse)

	assert.Truef(t, ok, fmt.Sprintf("Invalid type marshalled %t from response body, expected SuccessResponse type", object))
	assert.EqualValuesf(t, http.StatusInternalServerError, response.StatusCode, "Invalid http status code")
	assert.EqualValuesf(t, failResponse.Error, "Invalid user", "Invalid error message")
}

/*
	app.Get("/api/v1/users", object.getUsers) // DONE
	app.Get("/api/v1/users/email", object.getUser) // DONE
	app.Get(imagePath+":id", object.userImage().NewDownloader(), object.getImage)
	app.Post("/api/v1/users", object.userValidator().DuplicatedUser(), object.userImage().NewUploader(), object.newUser)
	app.Put("/api/v1/users", object.userImage().UpdateImage(), object.updateUser)
	app.Delete("/api/v1/users/email", object.userImage().DeleteImage(), object.deleteUser) //DONE

	https://stackoverflow.com/questions/43904974/testing-go-http-request-formfile
*/
