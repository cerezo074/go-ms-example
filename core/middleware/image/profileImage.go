package image

import (
	"io"
	"net/http"
	"strings"
	"user/app/utils/config"
	"user/app/utils/response"
	"user/core/dependencies/services"
	"user/core/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	PROFILE_IMAGE_USER_ENTITY     = "user_entity"
	PROFILE_IMAGE_FIELD           = "image_data"
	PROFILE_IMAGE__UPLOADED_ID    = "image_id"
	PROFILE_IMAGE_DOWNLOADED_FILE = "image_file"
)

type ProfileImageProvider struct {
	services.ProfileImageServices
	UserStore     entities.UserRepository
	UserValidator services.UserValidatorServices
	Credentials   config.Credentials
	Loader        services.ImageStorageLoader
}

func (object ProfileImageProvider) NewUploader() fiber.Handler {
	return func(context *fiber.Ctx) error {
		imageReader, err := getImageReader(context)
		if err != nil {
			//TODO: Log this error
			return context.Next()
		}

		fileID := uuid.New()
		filename := fileID.String()
		storageSession, err := object.Loader.Load(object.Credentials)
		if err != nil {
			//TODO: Log this error
			return context.Next()
		}

		imageURI, err := storageSession.Upload(imageReader, filename)
		context.Locals(PROFILE_IMAGE__UPLOADED_ID, imageURI)

		return context.Next()
	}
}

func (object ProfileImageProvider) NewDownloader() fiber.Handler {
	return func(context *fiber.Ctx) error {
		storageSession, err := object.Loader.Load(object.Credentials)
		if err != nil {
			return response.MakeErrorJSON(http.StatusInternalServerError, err.Error())
		}

		imageIDParam := context.Params(IMAGE_ID_KEY)
		result, err := storageSession.Download(imageIDParam)
		if err != nil {
			return response.MakeErrorJSON(http.StatusBadRequest, err.Error())
		}

		context.Locals(PROFILE_IMAGE_DOWNLOADED_FILE, result)

		return context.Next()
	}
}

func (object ProfileImageProvider) DeleteImage() fiber.Handler {
	return func(context *fiber.Ctx) error {
		email := context.Query(ADDRESS_KEY)
		user, filename := object.getUser(email, context, object.UserStore)
		if user == nil {
			return response.MakeErrorJSON(http.StatusNotFound, INVALID_USER_ERROR)
		}

		context.Locals(PROFILE_IMAGE_USER_ENTITY, *user)
		if filename == nil {
			return context.Next()
		}

		storageSession, err := object.Loader.Load(object.Credentials)
		if err != nil {
			return response.MakeErrorJSON(http.StatusBadRequest, err.Error())
		}

		err = storageSession.Delete(*filename)
		if err != nil {
			return response.MakeErrorJSON(http.StatusBadRequest, err.Error())
		}

		return context.Next()
	}
}

func (object ProfileImageProvider) UpdateImage() fiber.Handler {
	return func(context *fiber.Ctx) error {
		imageReader, err := getImageReader(context)
		if err != nil {
			//TODO: Log this error
			context.Locals(PROFILE_IMAGE__UPLOADED_ID, DEFAULT_IMAGE)
			return context.Next()
		}

		email := context.FormValue(EMAIL_KEY, "")
		user, filename := object.getUser(email, context, object.UserStore)
		if user == nil {
			return response.MakeErrorJSON(http.StatusNotFound, INVALID_USER_ERROR)
		}

		if filename == nil {
			randomID := uuid.New().String()
			filename = &randomID
		}

		storageSession, err := object.Loader.Load(object.Credentials)
		if err != nil {
			//TODO: Log this error
			return context.Next()
		}

		imageURI, err := storageSession.Upload(imageReader, *filename)
		context.Locals(PROFILE_IMAGE__UPLOADED_ID, imageURI)

		return context.Next()
	}
}

func getImageReader(context *fiber.Ctx) (io.Reader, error) {
	file, err := context.FormFile(PROFILE_IMAGE_FIELD)
	if err != nil {
		return nil, err
	}

	fileHeader, err := file.Open()
	if err != nil {
		return nil, err
	}

	return fileHeader, err
}

func (object ProfileImageProvider) getUser(email string, context *fiber.Ctx, userStore entities.UserRepository) (*entities.User, *string) {
	if !object.UserValidator.IsValidEmailFormat(email) {
		return nil, nil
	}

	user, err := userStore.User(email)
	if err != nil || user.Email == "" {
		return nil, nil
	}

	componentPaths := strings.Split(user.ImageID, "/")
	if len(componentPaths) <= 0 {
		return nil, nil
	}

	lastIndex := len(componentPaths) - 1
	lastComponent := componentPaths[lastIndex]
	if lastComponent == DEFAULT_IMAGE {
		return &user, nil
	}

	return &user, &lastComponent
}
