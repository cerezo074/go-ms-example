package routers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"user/core/dependencies/services"
	. "user/core/entities"
	image "user/core/middleware/image"
	"user/core/middleware/validator"
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
		Exist: func(email string) bool {
			for _, user := range users {
				if user.Email == email {
					return true
				}
			}

			return false
		},
		Save: func(newUser *User) error {
			if newUser == nil {
				return errors.New("Invalid new user to be saved, nil reference")
			}

			return nil
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
	deleteUserImage = FakeProfileImage{
		Delete: func(context *fiber.Ctx) error {
			context.Locals(image.PROFILE_IMAGE_USER_ENTITY, user1)
			return context.Next()
		},
	}
	deleteInvalidUserImage = FakeProfileImage{
		Delete: func(context *fiber.Ctx) error {
			context.Locals(image.PROFILE_IMAGE_USER_ENTITY, User{})
			return context.Next()
		},
	}
	angiImagePath    = "../../utils/assets/mioanji.jpg"
	anjiImageStorage = FakeImageLoader{
		DownloadImage: func(fileID string) (*services.ImageBufferedFile, error) {
			if fileID != "123456" {
				return nil, errors.New("Invalid file id reference")
			}

			buffer, size, err := utils.GetFile(angiImagePath)
			if err != nil {
				return nil, err
			}

			if size == 0 {
				return nil, errors.New("Invalid content file")
			}

			return &services.ImageBufferedFile{
				Data: buffer.Bytes(),
				Size: size,
			}, nil
		},
	}
)

func Test_GetUsers_WhenTheyExistInRepository(t *testing.T) {
	t.Parallel()
	server := utils.NewServer(utils.NewAllUsersUnmarshaller)
	fakeRepo := allUsersRepo
	appServices := NewUserMockedServices(fakeRepo, FakeValidator{}, FakeProfileImage{})

	routers.NewUserRouter().Register(server.FiberApp, appServices)
	response, object, _ := server.Execute("GET", "/api/v1/users", false, nil, nil)
	userSlice, ok := object.([]models.User)

	assert.Truef(t, ok, "Invalid type marshalled from response")
	assert.EqualValuesf(t, http.StatusOK, response.StatusCode, "Invalid http status code")
	assert.Lenf(t, userSlice, len(users), "Invalid length for users, must be 2")
	assert.EqualValuesf(t, userSlice[0].Email, "user1@gmail.com", "Invalid email for user1")
	assert.EqualValuesf(t, userSlice[1].Email, "user2@gmail.com", "Invalid email for user2")
}

func Test_ShouldntGetUsers_WhenTheyDontExistInRepository(t *testing.T) {
	t.Parallel()
	server := utils.NewServer(utils.NewAllUsersUnmarshaller)
	appServices := NewUserMockedServices(FakeRepo{}, FakeValidator{}, FakeProfileImage{})

	routers.NewUserRouter().Register(server.FiberApp, appServices)
	response, object, _ := server.Execute("GET", "/api/v1/users", false, nil, nil)
	userSlice, ok := object.([]models.User)

	assert.Truef(t, ok, fmt.Sprintf("Invalid type marshalled %t from response body, expected []User type", object))
	assert.EqualValuesf(t, http.StatusOK, response.StatusCode, "Invalid http status code")
	assert.Lenf(t, userSlice, 0, "Invalid length for users, must be 0")
}

func Test_GetUserByEmail_WhenItExistsInRepository(t *testing.T) {
	t.Parallel()
	server := utils.NewServer(utils.NewFindUserUnmarshaller)
	fakeRepo := findUserRepo
	appServices := NewUserMockedServices(fakeRepo, FakeValidator{}, FakeProfileImage{})

	routers.NewUserRouter().Register(server.FiberApp, appServices)
	response, object, _ := server.Execute("GET", "/api/v1/users/email?address="+user1.Email, false, nil, nil)
	user, ok := object.(models.User)

	assert.Truef(t, ok, fmt.Sprintf("Invalid type marshalled %t from response body, expected User type", object))
	assert.EqualValuesf(t, http.StatusOK, response.StatusCode, "Invalid http status code")
	assert.EqualValuesf(t, user.Email, "user1@gmail.com", "Invalid email for user1")
	assert.EqualValuesf(t, user.Nickname, "CR7", "Invalid nickname for user1")
}

func Test_ShouldntGetUserByEmail_WhenItDoesntExistInRepository(t *testing.T) {
	t.Parallel()
	server := utils.NewServer(utils.NewFailUnmarshaller)
	fakeRepo := findUserRepo
	appServices := NewUserMockedServices(fakeRepo, FakeValidator{}, FakeProfileImage{})
	invalidEmail := "invaliduser@test.com"

	routers.NewUserRouter().Register(server.FiberApp, appServices)
	response, object, _ := server.Execute("GET", "/api/v1/users/email?address="+invalidEmail, false, nil, nil)
	failResponse, ok := object.(models.FailResponse)

	assert.Truef(t, ok, fmt.Sprintf("Invalid type marshalled %t from response body, expected FailResponse type", object))
	assert.EqualValuesf(t, http.StatusNotFound, response.StatusCode, "Invalid http status code")
	assert.EqualValuesf(t, failResponse.Error, fmt.Sprintf("User doesn't exist with %v email", invalidEmail), "Invalid error message")
}

func Test_DeleteUserByEmail_WhenItExistsInRepository(t *testing.T) {
	t.Parallel()
	server := utils.NewServer(utils.NewSuccessUnmarshaller)
	fakeRepo := deleteUserRepo
	fakeImage := deleteUserImage
	appServices := NewUserMockedServices(fakeRepo, FakeValidator{}, fakeImage)

	routers.NewUserRouter().Register(server.FiberApp, appServices)
	response, object, _ := server.Execute("DELETE", "/api/v1/users/email?address="+user1.Email, false, nil, nil)
	successResponse, ok := object.(models.SuccessResponse)

	assert.Truef(t, ok, fmt.Sprintf("Invalid type marshalled %t from response body, expected SuccessResponse type", object))
	assert.EqualValuesf(t, http.StatusOK, response.StatusCode, "Invalid http status code")
	assert.EqualValuesf(t, successResponse.Data, "user deleted successfully", "Invalid error message")
}

func Test_ShouldntDeleteUserByEmail_WhenItDoesntExistInRepository(t *testing.T) {
	t.Parallel()
	server := utils.NewServer(utils.NewFailUnmarshaller)
	fakeRepo := deleteUserRepo
	fakeImage := deleteInvalidUserImage
	appServices := NewUserMockedServices(fakeRepo, FakeValidator{}, fakeImage)
	invalidEmail := "invaliduser@test.com"

	routers.NewUserRouter().Register(server.FiberApp, appServices)
	response, object, _ := server.Execute("DELETE", "/api/v1/users/email?address="+invalidEmail, false, nil, nil)
	failResponse, ok := object.(models.FailResponse)

	assert.Truef(t, ok, fmt.Sprintf("Invalid type marshalled %t from response body, expected SuccessResponse type", object))
	assert.EqualValuesf(t, http.StatusInternalServerError, response.StatusCode, "Invalid http status code")
	assert.EqualValuesf(t, failResponse.Error, "Invalid user", "Invalid error message")
}

func Test_GetUserImage_WhenItExists(t *testing.T) {
	server := utils.NewServer(utils.NewSuccessUnmarshaller)
	fakeRepo := FakeRepo{}
	fakeValidator := validator.UserValidatorProvider{UserStore: fakeRepo}
	fakeImageProvider := NewImageProvider(fakeRepo, fakeValidator, anjiImageStorage)
	appServices := NewUserMockedServices(fakeRepo, fakeValidator, fakeImageProvider)

	routers.NewUserRouter().Register(server.FiberApp, appServices)
	response, _, _ := server.Execute("GET", "/api/v1/users/image/123456", true, nil, nil)
	assert.NotNilf(t, response.Body, "Nil response")
	assert.Equal(t, response.StatusCode, http.StatusOK)
	responseBuffer := bytes.NewBuffer([]byte{})
	_, err := io.Copy(responseBuffer, response.Body)
	assert.NoErrorf(t, err, "Invalid response file")
	matched, _ := utils.FilesMatch(responseBuffer, angiImagePath)
	assert.Truef(t, matched, "Files didn't match")
}

func Test_DoesntGetImage_WhenItDoesntExist(t *testing.T) {
	server := utils.NewServer(utils.NewFailUnmarshaller)
	fakeRepo := FakeRepo{}
	fakeValidator := validator.UserValidatorProvider{UserStore: fakeRepo}
	fakeImageProvider := NewImageProvider(fakeRepo, fakeValidator, anjiImageStorage)
	appServices := NewUserMockedServices(fakeRepo, fakeValidator, fakeImageProvider)

	routers.NewUserRouter().Register(server.FiberApp, appServices)
	_, object, _ := server.Execute("GET", "/api/v1/users/image/654321", false, nil, nil)
	failResponse, ok := object.(models.FailResponse)
	assert.Truef(t, ok, fmt.Sprintf("Invalid response type, expect fail response, recieved %t", object))
	assert.Equalf(t, failResponse.Error, "Invalid file id reference", "Invalid error message")
}
