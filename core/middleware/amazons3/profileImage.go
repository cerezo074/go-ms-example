package amazons3

import (
	"io"
	"log"
	"net/url"
	"user/app/utils/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	S3_IMAGE_FIELD        = "image_data"
	S3_UPLOADED_IMAGE_URI = "image_uri"
	S3_ACL_POLICY         = "public-read"
	S3_URI_SCHEME         = "s3"
)

type AWSS3Config struct {
	AccessKeyID  string
	SecretKey    string
	S3Region     string
	S3BucketName string
}

func BuildAWSS3Config(credentials config.Credentials) AWSS3Config {
	return AWSS3Config{
		AccessKeyID:  credentials.AWSAccessKeyID,
		SecretKey:    credentials.AWSSecretKey,
		S3Region:     credentials.AWSS3ProfileRegion,
		S3BucketName: credentials.AWSS3ProfileBucket,
	}
}

func New(credentials config.Credentials) fiber.Handler {
	return func(context *fiber.Ctx) error {
		S3Credentials := BuildAWSS3Config(credentials)

		session, err := BuildAWSSession(S3Credentials)
		if err != nil {
			log.Fatalln(err.Error())
			return nil
		}

		imageReader, err := getImageReader(context)

		if err != nil {
			log.Fatalln(err.Error())
			return nil
		}

		imageURI, err := uploadImage(session, S3Credentials, imageReader, context)
		context.Locals(S3_UPLOADED_IMAGE_URI, imageURI)

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

func uploadImage(session *session.Session, config AWSS3Config, imageReader io.Reader, context *fiber.Ctx) (string, error) {
	uploader := s3manager.NewUploader(session)
	fileID := uuid.New()
	filename := fileID.String()

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.S3BucketName),
		ACL:    aws.String(S3_ACL_POLICY),
		Key:    aws.String(filename),
		Body:   imageReader,
	})

	if err != nil {
		return "", err
	}

	return (buildS3ImageURI(filename, config.S3BucketName)), nil
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

	return url.String()
}
