package mocks

import (
	"errors"
	"io"
	"user/app/utils/config"
	. "user/core/dependencies/services"
)

var (
	fakeCredentias = config.Credentials{}
)

type FakeImageLoader struct {
	UploadImage   func(io.Reader, string) (string, error)
	DownloadImage func(string) (*ImageBufferedFile, error)
	DeleteImage   func(string) error
}

func (object FakeImageLoader) Load(credentials config.Credentials) (ImageStorageSession, error) {
	return FakeImageStorage{
		UploadImage:   object.UploadImage,
		DownloadImage: object.DownloadImage,
		DeleteImage:   object.DeleteImage,
	}, nil
}

type FakeImageStorage struct {
	UploadImage   func(io.Reader, string) (string, error)
	DownloadImage func(string) (*ImageBufferedFile, error)
	DeleteImage   func(string) error
}

func (object FakeImageStorage) Upload(imageReader io.Reader, filename string) (string, error) {
	if object.UploadImage == nil {
		return "", errors.New("Upload image function not provided")
	}

	return object.UploadImage(imageReader, filename)
}

func (object FakeImageStorage) Download(imageIDParam string) (*ImageBufferedFile, error) {
	if object.DownloadImage == nil {
		return nil, errors.New("Download image function not provided")
	}

	return object.DownloadImage(imageIDParam)
}

func (object FakeImageStorage) Delete(objectID string) error {
	if object.DeleteImage == nil {
		return errors.New("Delete image function not provided")
	}

	return object.DeleteImage(objectID)
}
