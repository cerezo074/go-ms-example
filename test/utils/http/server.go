package http

import (
	"user/app/utils/response"

	"github.com/gofiber/fiber/v2"
)

func NewServer(unmarshaller ResponseUnmarshaller) *FakeServer {
	return &FakeServer{
		FiberApp: fiber.New(fiber.Config{
			ErrorHandler: response.HandleJSONError,
		}),
		Unmarshaller: unmarshaller,
	}
}
