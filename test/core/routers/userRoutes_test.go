package routers

import (
	"net/http"
	"testing"
	. "user/core/entities"
	"user/core/routers"
	utils "user/test/utils/http"
	. "user/test/utils/mocks"

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
)

func Test_ShouldGetUsers_WhenExistInRepository(t *testing.T) {
	server := &utils.FakeServer{
		FiberApp: fiber.New(),
	}

	fakeRepo := FakeRepo{
		AllUsers: func() ([]User, error) {
			return users, nil
		},
	}

	appServices := NewUserMockedServices(fakeRepo, FakeValidator{}, FakeImage{})
	routers.NewUserRouter().Register(server.FiberApp, appServices)
	response, JSONSlice, _ := server.Execute("GET", "/api/v1/users", nil, "data")

	assert.EqualValuesf(t, http.StatusOK, response.StatusCode, "Invalid http status code")
	assert.Lenf(t, JSONSlice, len(users), "Invalid length for users, must be 2")
	user1JSON := JSONSlice[0]
	user2JSON := JSONSlice[1]
	assert.EqualValuesf(t, user1JSON.Email, "user1@gmail.com", "Invalid email for user1")
	assert.EqualValuesf(t, user2JSON.Email, "user2@gmail.com", "Invalid email for user2")
}

func Test_ShouldntGetUsers_WhenRepositoryIsEmpty(t *testing.T) {
	server := &utils.FakeServer{
		FiberApp: fiber.New(),
	}

	appServices := NewUserMockedServices(FakeRepo{}, FakeValidator{}, FakeImage{})
	routers.NewUserRouter().Register(server.FiberApp, appServices)
	response, JSONSlice, _ := server.Execute("GET", "/api/v1/users", nil, "data")
	assert.EqualValuesf(t, http.StatusOK, response.StatusCode, "Invalid http status code")
	assert.Lenf(t, JSONSlice, 0, "Invalid length for users, must be 0")
}

/*
	app.Get("/api/v1/users", object.getUsers) // DONE
	app.Get("/api/v1/users/email", object.getUser)
	app.Get(imagePath+":id", object.userImage().NewDownloader(), object.getImage)
	app.Post("/api/v1/users", object.userValidator().DuplicatedUser(), object.userImage().NewUploader(), object.newUser)
	app.Put("/api/v1/users", object.userImage().UpdateImage(), object.updateUser)
	app.Delete("/api/v1/users/email", object.userImage().DeleteImage(), object.deleteUser)
*/
