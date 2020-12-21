package response

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

type JSONResponseType int

const (
	Success JSONResponseType = iota
	Fail    JSONResponseType = iota
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

func MakeJSON(responseType JSONResponseType, output *interface{}, err *ResponseError, context *fiber.Ctx) error {
	switch responseType {
	case Success:
		return context.JSON(fiber.Map{
			"data": output,
		})
	case Fail:
		jsonData, jsonError := err.JSON()
		if jsonError != nil {
			return err
		}

		context.Set("Content-Type", "application/json")
		context.Status(err.StatusCode).Send(jsonData)
		return nil
	default:
		return nil
	}
}
