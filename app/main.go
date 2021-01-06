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
	appStore  entities.Repository
	appConfig config.Credentials
)

func loadAppConfig() {
	config, err := config.LoadCredentials(".")

	if err != nil {
		log.Fatal(err)
		return
	}

	appConfig = config
}

func loadAppDB() {
	store, err := postgres.NewStore(appConfig.DBSource, appConfig.DBDriver)

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
		ErrorHandler: response.HandleJSONError,
	})

	userRouter := routers.NewUserRouter()
	userRouter.Register(appFiber, appStore, appConfig)
	log.Fatal(appFiber.Listen(appConfig.ServerAddress))
}
