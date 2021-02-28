package services

import "github.com/gofiber/fiber/v2"

type ImageServices struct {
	UserProfileImage S3ProfileImageServices
}

type S3ProfileImageServices interface {
	NewUploader() fiber.Handler
	NewDownloader() fiber.Handler
	DeleteImage() fiber.Handler
	UpdateImage() fiber.Handler
}
