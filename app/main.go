package main

import (
	"log"
	"user/app/db/postgres"
	"user/app/utils/config"
	"user/app/utils/response"
	"user/core/entities"
	"user/core/routers"

	"github.com/gofiber/fiber/v2"
)

var (
	appStore       entities.Repository
	appCredentials config.Credentials
)

func loadAppConfig() {
	credentials, err := config.LoadCredentials(".")

	if err != nil {
		log.Fatal(err)
		return
	}

	appCredentials = credentials
}

func loadAppDB() {
	store, err := postgres.NewStore(appCredentials.DBSource, appCredentials.DBDriver)

	if err != nil {
		log.Fatal(err)
		return
	}

	appStore = store
}

func init() {
	loadAppConfig()
	loadAppDB()
}

func main() {
	if appStore == nil {
		log.Fatal("Can not run web server before database server is up and ready")
		return
	}

	appFiber := fiber.New(fiber.Config{
		ErrorHandler: response.MakeErrorJSON,
	})
	userRouter := routers.NewUserRouter()
	userRouter.Register(appFiber, appStore)
	log.Fatal(appFiber.Listen(appCredentials.ServerAddress))
}
