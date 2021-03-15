package routers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"user/core/dependencies/services"
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
	elpibeImagePath = "../../utils/assets/elpibe.jpg"
	repeatedEmail   = "user1@gmail.com"
	user4           = User{
		ID:          uuid.New(),
		Email:       "user1@gmail.com",
		Nickname:    "CR7",
		Password:    "123456",
		ImageID:     "profile1.png",
		CountryCode: "USA",
		Birthday:    "07/01/2000",
	}
	user5 = User{
		ID:          uuid.New(),
		Email:       "user2@gmail.com",
		Nickname:    "Messi",
		Password:    "654321",
		ImageID:     "profile2.png",
		CountryCode: "USA",
		Birthday:    "08/01/2000",
	}
	savedUsers = []User{
		user4,
		user5,
	}
	repitedUser = utils.UserForm{
		Email:       repeatedEmail,
		Nickname:    "El Pibe'",
		Password:    "123456",
		ImagePath:   &elpibeImagePath,
		CountryCode: "COL",
		Birthday:    "12/22/2020",
	}
	elPibe = utils.UserForm{
		Email:       "carlos@valderrama.com",
		Nickname:    "El Pibe'",
		Password:    "123456",
		ImagePath:   &elpibeImagePath,
		CountryCode: "COL",
		Birthday:    "12/22/2020",
	}
	elPibeWithoutImage = utils.UserForm{
		Email:       "carlos@valderrama.com",
		Nickname:    "El Pibe'",
		Password:    "123456",
		ImagePath:   nil,
		CountryCode: "COL",
		Birthday:    "12/22/2020",
	}
	elPibesNewUserRepo = FakeRepo{
		AllUsers: func() ([]User, error) {
			return savedUsers, nil
		},
		Exist: func(email string) bool {
			for _, user := range savedUsers {
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

			if newUser.Email != elPibe.Email && newUser.Password != elPibe.Password {
				return errors.New(fmt.Sprintf("Invalid new user to be saved, %v", newUser))
			}

			return nil
		},
	}
	elPibesNewImageStorage = FakeImageLoader{
		UploadImage: func(image io.Reader, filename string) (string, error) {
			areFilesEquals, err := utils.FilesMatch(image, *elPibe.ImagePath)
			if err != nil || !areFilesEquals {
				return "", errors.New(fmt.Sprintf("Invalid file to be uploaded %s", filename))
			}

			return "fake-path/for-fake-image/" + filename, nil
		},
	}
)

func TestSignupUser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sign up User Suite")
}

var _ = Describe("Sign up User", func() {
	var server *utils.FakeServer
	var appServices services.App
	var requestBody io.Reader
	var contentTypeValue string
	var requestHeaders http.Header

	Context("User filled in register form", func() {
		When("User doesn't exist in repository", func() {
			BeforeEach(func() {
				requestBody, contentTypeValue, _ = utils.MultipartFormBody(elPibe)
				requestHeaders = http.Header{"Content-Type": []string{contentTypeValue}}
				server = utils.NewServer(utils.NewSuccessUnmarshaller)
				fakeValidator := validator.UserValidatorProvider{
					UserStore: elPibesNewUserRepo,
				}
				fakeImageProvider := NewImageProvider(elPibesNewUserRepo, fakeValidator, elPibesNewImageStorage)
				appServices = NewUserMockedServices(elPibesNewUserRepo, fakeValidator, fakeImageProvider)
			})

			It("Should get user crated message successfully", func() {
				routers.NewUserRouter().Register(server.FiberApp, appServices)
				response, object, _ := server.Execute("POST", "/api/v1/users", false, requestHeaders, requestBody)
				jsonResponse, ok := object.(models.SuccessResponse)
				Expect(ok).To(Equal(true))
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(jsonResponse.Data).To(Equal("user was created successfully"))
			})
		})

		When("User exists in repository", func() {
			BeforeEach(func() {
				requestBody, contentTypeValue, _ = utils.MultipartFormBody(repitedUser)
				requestHeaders = http.Header{"Content-Type": []string{contentTypeValue}}
				server = utils.NewServer(utils.NewFailUnmarshaller)
				fakeValidator := validator.UserValidatorProvider{
					UserStore: elPibesNewUserRepo,
				}
				fakeImageProvider := NewImageProvider(elPibesNewUserRepo, fakeValidator, elPibesNewImageStorage)
				appServices = NewUserMockedServices(elPibesNewUserRepo, fakeValidator, fakeImageProvider)
			})

			It("Shouldn't get user crated message successfully", func() {
				routers.NewUserRouter().Register(server.FiberApp, appServices)
				response, object, _ := server.Execute("POST", "/api/v1/users", false, requestHeaders, requestBody)
				jsonResponse, ok := object.(models.FailResponse)
				Expect(ok).To(Equal(true))
				Expect(response.StatusCode).To(Equal(http.StatusConflict))
				Expect(jsonResponse.Error).To(Equal(fmt.Sprintf("a user with the following email(%s) exist", repeatedEmail)))
			})
		})
	})

	Context("User doesn't send image inside form", func() {
		When("User doesn't exist in repository", func() {
			BeforeEach(func() {
				requestBody, contentTypeValue, _ = utils.MultipartFormBody(elPibeWithoutImage)
				requestHeaders = http.Header{"Content-Type": []string{contentTypeValue}}
				server = utils.NewServer(utils.NewSuccessUnmarshaller)
				fakeValidator := validator.UserValidatorProvider{
					UserStore: elPibesNewUserRepo,
				}
				fakeImageProvider := NewImageProvider(elPibesNewUserRepo, fakeValidator, elPibesNewImageStorage)
				appServices = NewUserMockedServices(elPibesNewUserRepo, fakeValidator, fakeImageProvider)
			})

			It("Should get user crated message successfully", func() {
				routers.NewUserRouter().Register(server.FiberApp, appServices)
				response, object, _ := server.Execute("POST", "/api/v1/users", false, requestHeaders, requestBody)
				jsonResponse, ok := object.(models.SuccessResponse)
				Expect(ok).To(Equal(true))
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(jsonResponse.Data).To(Equal("user was created successfully"))
			})
		})
	})
})
