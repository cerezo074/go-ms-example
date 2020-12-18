package main

import (
	"user/core/routers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	userRouter := routers.BuildUserRouter()
	userRouter.Register(app)
	app.Listen(":3000")
}
