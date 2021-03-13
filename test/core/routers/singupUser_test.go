package routers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"

	"user/core/dependencies/services"
	"user/core/entities"
	"user/core/middleware/validator"
	"user/core/routers"
	utils "user/test/utils/http"
	. "user/test/utils/mocks"
	"user/test/utils/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	repeatedEmail = "user1@gmail.com"
	repitedUser   = UserForm{
		email:       repeatedEmail,
		nickname:    "El Pibe'",
		password:    "123456",
		imagePath:   "../../utils/assets/shishio.jpg",
		countryCode: "COL",
		birthday:    "12/22/2020",
	}
	elPibe = UserForm{
		email:       "carlos@valderrama.com",
		nickname:    "El Pibe'",
		password:    "123456",
		imagePath:   "../../utils/assets/shishio.jpg",
		countryCode: "COL",
		birthday:    "12/22/2020",
	}
	elPibesNewUserRepo = FakeRepo{
		AllUsers: func() ([]entities.User, error) {
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
		Save: func(newUser *entities.User) error {
			if newUser == nil {
				return errors.New("Invalid new user to be saved, nil reference")
			}

			if newUser.Email != elPibe.email && newUser.Password != elPibe.password {
				return errors.New(fmt.Sprintf("Invalid new user to be saved, %v", newUser))
			}

			return nil
		},
	}
	elPibesNewImageStorage = FakeImageLoader{
		UploadImage: func(image io.Reader, filename string) (string, error) {
			areFilesEquals, err := FilesMatch(image, elPibe.imagePath)
			if err != nil || !areFilesEquals {
				return "", errors.New(fmt.Sprintf("Invalid file to be uploaded %s", filename))
			}

			return "fake-path/for-fake-image/" + filename, nil
		},
	}
)

func TestCart(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Signing up Suite")
}

var _ = Describe("Signing up", func() {
	var server *utils.FakeServer
	var appServices services.App
	var requestBody io.Reader
	var contentTypeValue string
	var requestHeaders http.Header

	Context("User filled in register form", func() {
		// When("User doesn't exist in repository", func() {
		// 	BeforeEach(func() {
		// 		requestBody, contentTypeValue, _ = MultipartFormBody(elPibe)
		// 		requestHeaders = http.Header{"Content-Type": []string{contentTypeValue}}
		// 		server = buildServer(utils.NewSuccessUnmarshaller)
		// 		fakeValidator := validator.UserValidatorProvider{
		// 			UserStore: elPibesNewUserRepo,
		// 		}
		// 		fakeImageProvider := NewImageProvider(elPibesNewUserRepo, fakeValidator, elPibesNewImageStorage)
		// 		appServices = NewUserMockedServices(elPibesNewUserRepo, fakeValidator, fakeImageProvider)
		// 	})

		// 	It("Should get user crated message successfully", func() {
		// 		routers.NewUserRouter().Register(server.FiberApp, appServices)
		// 		response, object, _ := server.Execute("POST", "/api/v1/users", requestHeaders, requestBody)
		// 		jsonResponse, ok := object.(models.SuccessResponse)
		// 		Expect(ok).To(Equal(true))
		// 		Expect(response.StatusCode).To(Equal(http.StatusOK))
		// 		Expect(jsonResponse.Data).To(Equal("user was created successfully"))
		// 	})
		// })

		When("User exists in repository", func() {
			BeforeEach(func() {
				requestBody, contentTypeValue, _ = MultipartFormBody(repitedUser)
				requestHeaders = http.Header{"Content-Type": []string{contentTypeValue}}
				server = buildServer(utils.NewFailUnmarshaller)
				fakeValidator := validator.UserValidatorProvider{
					UserStore: elPibesNewUserRepo,
				}
				fakeImageProvider := NewImageProvider(elPibesNewUserRepo, fakeValidator, elPibesNewImageStorage)
				appServices = NewUserMockedServices(elPibesNewUserRepo, fakeValidator, fakeImageProvider)
			})

			It("Shouldn't get user crated message successfully", func() {
				routers.NewUserRouter().Register(server.FiberApp, appServices)
				response, object, _ := server.Execute("POST", "/api/v1/users", requestHeaders, requestBody)
				jsonResponse, ok := object.(models.FailResponse)
				Expect(ok).To(Equal(true))
				Expect(response.StatusCode).To(Equal(http.StatusConflict))
				Expect(jsonResponse.Error).To(Equal(fmt.Sprintf("a user with the following email(%s) exist", repeatedEmail)))
			})
		})
	})
})

func FilesMatch(rawFile io.Reader, rigthFilePath string) (bool, error) {
	rawImage, err := os.Open(elPibe.imagePath)
	if err != nil {
		return false, err
	}

	fileInfo, err := rawImage.Stat()
	if err != nil {
		return false, err
	}

	bytes := bytes.NewBuffer([]byte{})
	bytesWrited, err := io.Copy(bytes, rawFile)
	if err != nil {
		return false, err
	}

	if fileInfo.Size() != bytesWrited {
		return false, errors.New("Files dont contain same size")
	}

	return true, nil
}

type UserForm struct {
	email       string `form:"email"`
	nickname    string `form:"nickname"`
	password    string `form:"password"`
	imagePath   string `form:"image_data" type:"file"`
	countryCode string `form:"country_code"`
	birthday    string `form:"birthday"`
}

func (object UserForm) imageName() string {
	_, filename := path.Split(object.imagePath)
	return filename
}

func AddMultipartFile(key string, filepath string, writer *multipart.Writer) error {
	_, filename := path.Split(filepath)
	if len(filename) == 0 || strings.Contains(filename, " ") {
		return errors.New(fmt.Sprintf("Invalid filename %v", filename))
	}

	part, err := writer.CreateFormFile(key, filename)
	if err != nil {
		return err
	}

	sample, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer sample.Close()
	_, err = io.Copy(part, sample)
	if err != nil {
		return err
	}

	err = writer.Close()
	return nil
}

func AddMultipartField(key string, value string, writer *multipart.Writer) error {
	field, err := writer.CreateFormField(key)
	if err != nil {
		return err
	}

	_, err = field.Write([]byte(value))
	return nil
}

// func multipartFormBody(form interface{}) (*bytes.Buffer, string, error) {
// 	body := new(bytes.Buffer)
// 	writer := multipart.NewWriter(body)
// 	value := reflect.ValueOf(form).Elem()

// 	for i := 0; i < value.NumField(); i++ {
// 		typeField := value.Type().Field(i)
// 		valueField := value.Field(i)
// 		tag := typeField.Tag

// 		formKey := tag.Get("form")
// 		if formKey == "" {
// 			continue
// 		}

// 		formValue := valueField.Interface()
// 		stringValue, ok := formValue.(string)
// 		if !ok {
// 			continue
// 		}

// 		if tag.Get("type") == "file" {
// 			err := AddMultipartFile(formKey, stringValue, writer)
// 			if err != nil {
// 				return nil, "", err
// 			}
// 		} else {
// 			err := AddMultipartField(formKey, stringValue, writer)
// 			if err != nil {
// 				return nil, "", err
// 			}
// 		}
// 	}

// 	return body, writer.FormDataContentType(), nil
// }

func MultipartFormBody(form UserForm) (*bytes.Buffer, string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	err := AddMultipartField("email", form.email, writer)
	if err != nil {
		return nil, "", err
	}

	err = AddMultipartField("nickname", form.nickname, writer)
	if err != nil {
		return nil, "", err
	}

	err = AddMultipartField("password", form.password, writer)
	if err != nil {
		return nil, "", err
	}

	err = AddMultipartField("country_code", form.countryCode, writer)
	if err != nil {
		return nil, "", err
	}

	err = AddMultipartField("birthday", form.birthday, writer)
	if err != nil {
		return nil, "", err
	}

	err = AddMultipartFile("image_data", form.imagePath, writer)
	if err != nil {
		return nil, "", err
	}

	writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}
