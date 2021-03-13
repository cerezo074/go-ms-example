package services

import (
	"io"
	"user/app/utils/config"

	"github.com/gofiber/fiber/v2"
)

type ImageServices struct {
	UserProfileImage ProfileImageServices
}

type ProfileImageServices interface {
	NewUploader() fiber.Handler
	NewDownloader() fiber.Handler
	DeleteImage() fiber.Handler
	UpdateImage() fiber.Handler
}

type ImageStorageLoader interface {
	Load(credentials config.Credentials) (ImageStorageSession, error)
}

type ImageStorageSession interface {
	Upload(imageReader io.Reader, filename string) (string, error)
	Download(imageIDParam string) (*ImageBufferedFile, error)
	Delete(objectID string) error
}

type ImageBufferedFile struct {
	Data []byte
	Size int64
}
