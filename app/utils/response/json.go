package response

import (
	"github.com/gofiber/fiber/v2"
)

type JSONResponseType int

const (
	Success JSONResponseType = iota
	Fail    JSONResponseType = iota
)

func MakeJSON(responseType JSONResponseType, output interface{}, err error, context *fiber.Ctx) error {
	switch responseType {
	case Success:
		return context.JSON(fiber.Map{
			"data": output,
		})
	case Fail:
		return err
	default:
		return nil
	}
}
