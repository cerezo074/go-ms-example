package validator

import (
	"net/http"
	"regexp"
	"user/app/utils/response"
	"user/core/entities"

	"github.com/gofiber/fiber/v2"
)

const (
	EMAIL_FIELD         = "email"
	EMAIL_REGEX_PATTERN = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

func DuplicatedUser(userStore entities.UserRepository) fiber.Handler {
	return func(context *fiber.Ctx) error {
		email := context.FormValue(EMAIL_FIELD)

		if !IsValidEmailFormat(email) {
			return buildInvalidFormatError(email)
		}

		if userStore.ExistUser(email) {
			return buildDuplicatedError(email)
		}

		return context.Next()
	}
}

func buildDuplicatedError(email string) error {
	return response.ResponseError{StatusCode: http.StatusConflict, Message: "a user with the following email(" + email + ") exist"}
}

func buildInvalidFormatError(email string) error {
	return response.ResponseError{StatusCode: http.StatusBadRequest, Message: "invalid email(" + email + ") format"}
}

func IsValidEmailFormat(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}

	var emailRegex = regexp.MustCompile(EMAIL_REGEX_PATTERN)
	return emailRegex.MatchString(email)
}
