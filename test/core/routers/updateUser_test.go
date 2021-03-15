package routers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"user/core/dependencies/services"
	"user/core/entities"
	. "user/core/entities"
	"user/core/middleware/validator"
	"user/core/routers"
	utils "user/test/utils/http"
	. "user/test/utils/mocks"
	"user/test/utils/models"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	makotoShishio = User{
		ID:          uuid.New(),
		Email:       "makoto@shishio.com",
		Nickname:    "Mummy",
		Password:    "555555",
		ImageID:     "profile20.png",
		CountryCode: "JPN",
		Birthday:    "07/01/1900",
	}
	user3 = User{
		ID:          uuid.New(),
		Email:       "user3@gmail.com",
		Nickname:    "Maldini",
		Password:    "444444",
		ImageID:     "profile3.png",
		CountryCode: "ITA",
		Birthday:    "08/01/2020",
	}
	shishioImageUpdated = utils.UserForm{
		Email:       "makoto@shishio.com",
		Nickname:    "Mummy",
		Password:    "333333",
		ImagePath:   &shishioImagePath,
		CountryCode: "UK",
		Birthday:    "12/22/2020",
	}
	shishioWithoutImage = utils.UserForm{
		Email:       "makoto@shishio.com",
		Nickname:    "Mummy",
		Password:    "111111",
		ImagePath:   nil,
		CountryCode: "UK",
		Birthday:    "12/22/2020",
	}
	invalidUpdatedUser = utils.UserForm{
		Email:       "newUser@test.com",
		Nickname:    "A New User",
		Password:    "222222",
		ImagePath:   &shishioImagePath,
		CountryCode: "COL",
		Birthday:    "12/22/2020",
	}
	updatedUsers          = []User{makotoShishio, user3}
	shishioImagePath      = "../../utils/assets/shishio.jpg"
	searchAllUsersByEmail = func(email string) (entities.User, error) {
		for _, user := range updatedUsers {
			if user.Email == email {
				return user, nil
			}
		}

		return entities.User{}, errors.New("Invalid user to be searched by email")
	}
	updateShisioUserOnRepo = func(updatedUser *entities.User) error {
		if updatedUser == nil {
			return errors.New("Invalid new user to be saved, nil reference")
		}

		if updatedUser.Email != makotoShishio.Email || updatedUser.Password == makotoShishio.Password {
			return errors.New(fmt.Sprintf("Invalid old user to be updated, %v", updatedUser))
		}

		return nil
	}
	updateShishioRepo = FakeRepo{
		UserByEmail: searchAllUsersByEmail,
		Update:      updateShisioUserOnRepo,
	}
	updateShishioImageStorage = FakeImageLoader{
		UploadImage: func(image io.Reader, filename string) (string, error) {
			areFilesEquals, err := utils.FilesMatch(image, *shishioImageUpdated.ImagePath)
			if err != nil || !areFilesEquals {
				return "", errors.New(fmt.Sprintf("Invalid file to be uploaded %s", filename))
			}

			return "new-path/for-updated-image/" + "updatedImage.png", nil
		},
	}
)

func TestUpdateUser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Update User Suite")
}

var _ = Describe("Update User", func() {
	var server *utils.FakeServer
	var appServices services.App
	var requestBody io.Reader
	var contentTypeValue string
	var requestHeaders http.Header

	Context("User updates all fields in form except email", func() {
		When("User exists in repository", func() {
			BeforeEach(func() {
				requestBody, contentTypeValue, _ = utils.MultipartFormBody(shishioImageUpdated)
				requestHeaders = http.Header{"Content-Type": []string{contentTypeValue}}
				server = utils.NewServer(utils.NewSuccessUnmarshaller)
				fakeValidator := validator.UserValidatorProvider{
					UserStore: updateShishioRepo,
				}
				fakeImageProvider := NewImageProvider(updateShishioRepo, fakeValidator, updateShishioImageStorage)
				appServices = NewUserMockedServices(updateShishioRepo, fakeValidator, fakeImageProvider)
			})

			It("Should get user crated message successfully", func() {
				routers.NewUserRouter().Register(server.FiberApp, appServices)
				response, object, _ := server.Execute("PUT", "/api/v1/users", false, requestHeaders, requestBody)
				jsonResponse, ok := object.(models.SuccessResponse)
				Expect(ok).To(Equal(true))
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(jsonResponse.Data).To(Equal("user updated successfully"))
			})
		})

		When("User doesn't exist in repository", func() {
			BeforeEach(func() {
				requestBody, contentTypeValue, _ = utils.MultipartFormBody(invalidUpdatedUser)
				requestHeaders = http.Header{"Content-Type": []string{contentTypeValue}}
				server = utils.NewServer(utils.NewFailUnmarshaller)
				fakeValidator := validator.UserValidatorProvider{
					UserStore: updateShishioRepo,
				}
				fakeImageProvider := NewImageProvider(updateShishioRepo, fakeValidator, updateShishioImageStorage)
				appServices = NewUserMockedServices(updateShishioRepo, fakeValidator, fakeImageProvider)
			})

			It("Should get error message when attempt to create user crated message successfully", func() {
				routers.NewUserRouter().Register(server.FiberApp, appServices)
				response, object, _ := server.Execute("PUT", "/api/v1/users", false, requestHeaders, requestBody)
				jsonResponse, ok := object.(models.FailResponse)
				Expect(ok).To(Equal(true))
				Expect(response.StatusCode).To(Equal(http.StatusNotFound))
				Expect(jsonResponse.Error).To(Equal("invalid user"))
			})
		})
	})

	Context("User doesn't send an image inside form", func() {
		When("User exists in repository", func() {
			BeforeEach(func() {
				requestBody, contentTypeValue, _ = utils.MultipartFormBody(shishioWithoutImage)
				requestHeaders = http.Header{"Content-Type": []string{contentTypeValue}}
				server = utils.NewServer(utils.NewSuccessUnmarshaller)
				fakeValidator := validator.UserValidatorProvider{
					UserStore: updateShishioRepo,
				}
				fakeImageProvider := NewImageProvider(updateShishioRepo, fakeValidator, updateShishioImageStorage)
				appServices = NewUserMockedServices(updateShishioRepo, fakeValidator, fakeImageProvider)
			})

			It("Should get user updated message successfully", func() {
				routers.NewUserRouter().Register(server.FiberApp, appServices)
				response, object, _ := server.Execute("PUT", "/api/v1/users", false, requestHeaders, requestBody)
				jsonResponse, ok := object.(models.SuccessResponse)
				Expect(ok).To(Equal(true))
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(jsonResponse.Data).To(Equal("user updated successfully"))
			})
		})
	})
})
