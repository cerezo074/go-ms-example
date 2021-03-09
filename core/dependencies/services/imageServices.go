package services

import "github.com/gofiber/fiber/v2"

type ImageServices struct {
	UserProfileImage ProfileImageServices
}

type ProfileImageServices interface {
	NewUploader() fiber.Handler
	NewDownloader() fiber.Handler
	DeleteImage() fiber.Handler
	UpdateImage() fiber.Handler
}
