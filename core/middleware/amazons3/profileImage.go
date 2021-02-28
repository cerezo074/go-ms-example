package amazons3

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"user/app/utils/config"
	"user/app/utils/response"
	"user/core/dependencies/services"
	"user/core/entities"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	S3_USER_ENTITY           = "user_entity"
	S3_IMAGE_FIELD           = "image_data"
	S3_UPLOADED_IMAGE_ID     = "image_id"
	S3_DOWNLOADED_IMAGE_FILE = "image_file"
	S3_ACL_POLICY            = ""
	S3_URI_SCHEME            = "s3"
	EMAIL_KEY                = "email"
	ADDRESS_KEY              = "address"
	IMAGE_ID_KEY             = "id"
	INVALID_USER_ERROR       = "Invalid User"
	DEFAULT_IMAGE            = "default_1.jpg"
)

type S3ProfileImageProvider struct {
	services.S3ProfileImageServices
	Credentials   config.Credentials
	UserStore     entities.UserRepository
	UserValidator services.UserValidatorServices
}

type AWSS3Config struct {
	AccessKeyID  string
	SecretKey    string
	S3Region     string
	S3BucketName string
}

type AWSBufferedFile struct {
	Data []byte
	Size int64
}

func BuildAWSS3Config(credentials config.Credentials) AWSS3Config {
	return AWSS3Config{
		AccessKeyID:  credentials.AWSAccessKeyID,
		SecretKey:    credentials.AWSSecretKey,
		S3Region:     credentials.AWSS3ProfileRegion,
		S3BucketName: credentials.AWSS3ProfileBucket,
	}
}

func (object S3ProfileImageProvider) NewUploader() fiber.Handler {
	return func(context *fiber.Ctx) error {
		imageReader, err := getImageReader(context)
		if err != nil {
			//TODO: Log this error
			return context.Next()
		}

		fileID := uuid.New()
		filename := fileID.String()
		S3Credentials := BuildAWSS3Config(object.Credentials)
		session, err := BuildAWSSession(S3Credentials)
		if err != nil {
			//TODO: Log this error
			return context.Next()
		}

		imageURI, err := uploadImage(session, S3Credentials, imageReader, filename, context)
		context.Locals(S3_UPLOADED_IMAGE_ID, imageURI)

		return context.Next()
	}
}

func (object S3ProfileImageProvider) NewDownloader() fiber.Handler {
	return func(context *fiber.Ctx) error {
		S3Credentials := BuildAWSS3Config(object.Credentials)
		session, err := BuildAWSSession(S3Credentials)
		if err != nil {
			return response.MakeErrorJSON(http.StatusInternalServerError, err.Error())
		}

		result, err := downloadImage(session, S3Credentials, context)
		if err != nil {
			return err
		}

		context.Locals(S3_DOWNLOADED_IMAGE_FILE, result)

		return context.Next()
	}
}

func (object S3ProfileImageProvider) DeleteImage() fiber.Handler {
	return func(context *fiber.Ctx) error {
		email := context.Query(ADDRESS_KEY)
		user, filename := object.getUser(email, context, object.UserStore)
		if user == nil {
			return response.MakeErrorJSON(http.StatusNotFound, INVALID_USER_ERROR)
		}

		context.Locals(S3_USER_ENTITY, *user)
		if filename == nil {
			return context.Next()
		}

		S3Credentials := BuildAWSS3Config(object.Credentials)
		session, err := BuildAWSSession(S3Credentials)
		if err != nil {
			return response.MakeErrorJSON(http.StatusBadRequest, err.Error())
		}

		err = deleteProfileImage(session, S3Credentials, *filename)
		if err != nil {
			return response.MakeErrorJSON(http.StatusBadRequest, err.Error())
		}

		return context.Next()
	}
}

func (object S3ProfileImageProvider) UpdateImage() fiber.Handler {
	return func(context *fiber.Ctx) error {
		imageReader, err := getImageReader(context)
		if err != nil {
			//TODO: Log this error
			context.Locals(S3_UPLOADED_IMAGE_ID, DEFAULT_IMAGE)
			return context.Next()
		}

		email := context.FormValue(EMAIL_KEY, "")
		user, filename := object.getUser(email, context, object.UserStore)
		if user == nil {
			return response.MakeErrorJSON(http.StatusNotFound, INVALID_USER_ERROR)
		}

		if filename == nil {
			randomID := uuid.New().String()
			filename = &randomID
		}

		S3Credentials := BuildAWSS3Config(object.Credentials)
		session, err := BuildAWSSession(S3Credentials)
		if err != nil {
			//TODO: Log this error
			return context.Next()
		}

		imageURI, err := uploadImage(session, S3Credentials, imageReader, *filename, context)
		context.Locals(S3_UPLOADED_IMAGE_ID, imageURI)

		return context.Next()
	}
}

func BuildAWSSession(config AWSS3Config) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(config.S3Region),
		Credentials: credentials.NewStaticCredentials(
			config.AccessKeyID,
			config.SecretKey,
			"",
		),
	})
}

func downloadImage(session *session.Session, config AWSS3Config, context *fiber.Ctx) (*AWSBufferedFile, error) {
	imageIDParam := context.Params(IMAGE_ID_KEY)
	imageID := strings.ReplaceAll(imageIDParam, " ", "")
	if len(imageID) <= 0 {
		imageID = DEFAULT_IMAGE
	}

	downloader := s3manager.NewDownloader(session)
	s3object := &s3.GetObjectInput{
		Bucket: aws.String(config.S3BucketName),
		Key:    aws.String(imageID),
	}

	myFile := &aws.WriteAtBuffer{}
	bytesDownloaded, err := downloader.Download(myFile, s3object)
	if err != nil {
		return nil, response.MakeErrorJSON(http.StatusBadGateway, err.Error())
	}

	data := myFile.Bytes()

	return &AWSBufferedFile{Data: data, Size: bytesDownloaded}, nil
}

func uploadImage(session *session.Session, config AWSS3Config, imageReader io.Reader, fileID string, context *fiber.Ctx) (string, error) {
	uploader := s3manager.NewUploader(session)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.S3BucketName),
		ACL:    aws.String(S3_ACL_POLICY),
		Key:    aws.String(fileID),
		Body:   imageReader,
	})

	if err != nil {
		return "", err
	}

	return (buildS3ImageURI(fileID, config.S3BucketName)), nil
}

func getImageReader(context *fiber.Ctx) (io.Reader, error) {
	file, err := context.FormFile(S3_IMAGE_FIELD)
	if err != nil {
		return nil, err
	}

	fileHeader, err := file.Open()
	if err != nil {
		return nil, err
	}

	return fileHeader, err
}

func buildS3ImageURI(imageID string, bucketName string) string {
	url := url.URL{
		Scheme: S3_URI_SCHEME,
		Host:   bucketName,
		Path:   imageID,
	}

	return url.Path
}

func deleteProfileImage(session *session.Session, config AWSS3Config, objectID string) error {
	svc := s3.New(session)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(config.S3BucketName), Key: aws.String(objectID)})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(config.S3BucketName),
		Key:    aws.String(objectID),
	})

	return err
}

func (object S3ProfileImageProvider) getUser(email string, context *fiber.Ctx, userStore entities.UserRepository) (*entities.User, *string) {
	if !object.UserValidator.IsValidEmailFormat(email) {
		return nil, nil
	}

	user, err := userStore.User(email)
	if err != nil {
		return nil, nil
	}

	componentPaths := strings.Split(user.ImageID, "/")
	if len(componentPaths) <= 0 {
		return nil, nil
	}

	lastIndex := len(componentPaths) - 1
	lastComponent := componentPaths[lastIndex]
	if lastComponent == DEFAULT_IMAGE {
		return &user, nil
	}

	return &user, &lastComponent
}
