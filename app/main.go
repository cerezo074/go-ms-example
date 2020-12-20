package main

import (
	"log"
	"user/app/db/postgres"
	"user/core/entities"
	"user/core/routers"

	"github.com/gofiber/fiber/v2"
)

var (
	store entities.Repository
)

func init() {
	postgresStore, err := postgres.NewStore("postgres://admin:password@localhost:5432/freefortalking?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	} else {
		store = postgresStore
	}
}

func main() {
	if store == nil {
		log.Fatal("Can not run web server before database server is up and ready")
		return
	}

	app := fiber.New()
	userRouter := routers.NewUserRouter()
	userRouter.Register(app, store)
	log.Fatal(app.Listen(":3000"))
}
