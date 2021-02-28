package validator

import (
	"errors"
	"net/http"
	"regexp"
	"user/app/utils/response"
	"user/core/entities"
	"user/core/services"

	"github.com/gofiber/fiber/v2"
)

const (
	EMAIL_FIELD         = "email"
	EMAIL_REGEX_PATTERN = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

type UserValidatorProvider struct {
	services.UserValidatorServices
	UserStore entities.UserRepository
}

func (object UserValidatorProvider) DuplicatedUser() fiber.Handler {
	return func(context *fiber.Ctx) error {
		email := context.FormValue(EMAIL_FIELD)

		if !object.IsValidEmailFormat(email) {
			return buildInvalidFormatError(email)
		}

		if object.UserStore.ExistUser(email) {
			return buildDuplicatedError(email)
		}

		return context.Next()
	}
}

func (object UserValidatorProvider) IsValidEmailFormat(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}

	var emailRegex = regexp.MustCompile(EMAIL_REGEX_PATTERN)
	return emailRegex.MatchString(email)
}

func (object UserValidatorProvider) IsValid(user entities.User) error {
	errorMessage := ""
	if user.ID.String() == "" {
		errorMessage += "Invalid id, "
	}

	if user.Email == "" {
		errorMessage += "Invalid email, "
	}

	if user.Nickname == "" {
		errorMessage += "Invalid nickname, "
	}

	if user.Password == "" {
		errorMessage += "Invalid password, "
	}

	if user.ImageID == "" {
		errorMessage += "Invalid image id, "
	}

	if user.CountryCode == "" {
		errorMessage += "Invalid country code, "
	}

	if user.Birthday == "" {
		errorMessage += "Invalid birthday"
	}

	if errorMessage == "" {
		return nil
	}

	return errors.New("There are some fields invalids: " + errorMessage)
}

func buildDuplicatedError(email string) error {
	return response.ResponseError{StatusCode: http.StatusConflict, Message: "a user with the following email(" + email + ") exist"}
}

func buildInvalidFormatError(email string) error {
	return response.ResponseError{StatusCode: http.StatusBadRequest, Message: "invalid email(" + email + ") format"}
}
