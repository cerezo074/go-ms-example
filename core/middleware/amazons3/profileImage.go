package amazons3

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"user/app/utils/config"
	"user/app/utils/response"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	S3_IMAGE_FIELD        = "image_data_field"
	S3_UPLOADED_IMAGE_URI = "image_uri"
	S3_DOWNLOADED_IMAGE   = "image_file"
	S3_ACL_POLICY         = ""
	S3_URI_SCHEME         = "s3"
	DEFAULT_IMAGE         = "133702635_2921149807988560_5555061904179489233_o.jpg"
)

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

func NewUploader(credentials config.Credentials) fiber.Handler {
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

func NewDownloader(credentials config.Credentials) fiber.Handler {
	return func(context *fiber.Ctx) error {
		S3Credentials := BuildAWSS3Config(credentials)

		session, err := BuildAWSSession(S3Credentials)
		if err != nil {
			log.Fatalln(err.Error())
			return nil
		}

		result, err := downloadImage(session, S3Credentials, context)

		if err != nil {
			return err
		}

		context.Locals(S3_DOWNLOADED_IMAGE, result)

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
	imageIDParam := context.Params("id")
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
	data := myFile.Bytes()

	if err != nil {
		return nil, response.MakeErrorJSON(http.StatusBadGateway, err.Error())
	}

	return &AWSBufferedFile{Data: data, Size: bytesDownloaded}, nil
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

	return url.Path
}
