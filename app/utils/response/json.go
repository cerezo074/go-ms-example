package response

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

const (
	data = "data"
)

type JSONError struct {
	Reason string `json:"error"`
}

type ResponseError struct {
	StatusCode int
	Message    string
}

func (e ResponseError) Error() string {
	return e.Message
}

func (e ResponseError) JSON() ([]byte, error) {
	data := JSONError{Reason: e.Message}
	encodedData, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	return encodedData, nil
}

func MakeSuccessJSON(output interface{}, context *fiber.Ctx) error {
	return context.JSON(fiber.Map{
		data: output,
	})
}

func HandleJSONError(context *fiber.Ctx, err error) error {
	if responseError, ok := err.(ResponseError); ok {
		jsonData, jsonError := responseError.JSON()
		if jsonError != nil {
			return err
		}

		context.Set("Content-Type", "application/json")
		return context.Status(responseError.StatusCode).Send(jsonData)
	}

	if err != nil {
		return context.Status(500).SendString("Internal Server Error")
	}

	return nil
}

func MakeErrorJSON(httpStatusCode int, description string) error {
	return sendError(httpStatusCode, description)
}

func sendError(httpStatusCode int, description string) error {
	return ResponseError{
		StatusCode: httpStatusCode,
		Message:    description,
	}
}
