package main

import (
	"log"
	"user/app/utils/response"
	"user/core/dependency"
	"user/core/routers"
	"user/core/services"

	"github.com/gofiber/fiber/v2"
)

type server struct {
	app           *fiber.App
	serverAddress string
}

func (object *server) start() {
	log.Fatal(object.app.Listen(object.serverAddress))
}

func SetupApp(dependencies *services.App) (*server, error) {
	appDependencies := dependencies
	if appDependencies == nil {
		defaultDependencies, err := dependency.NewServiceLocator(nil)
		if err != nil {
			return nil, err
		}

		appDependencies = defaultDependencies
	}

	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: response.HandleJSONError,
	})

	userRouter := routers.NewUserRouter()
	userRouter.Register(fiberApp, *appDependencies)

	return &server{
		app:           fiberApp,
		serverAddress: appDependencies.Credentials.ServerAddress,
	}, nil
}

func main() {
	server, err := SetupApp(nil)
	if err != nil {
		log.Fatalf("Can't init app, %v", err)
		return
	}

	server.start()
}
