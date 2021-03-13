package image

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"user/app/utils/config"
	"user/app/utils/response"
	"user/core/dependencies/services"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	S3_ACL_POLICY      = ""
	S3_URI_SCHEME      = "s3"
	EMAIL_KEY          = "email"
	ADDRESS_KEY        = "address"
	IMAGE_ID_KEY       = "id"
	INVALID_USER_ERROR = "invalid user"
	DEFAULT_IMAGE      = "default_1.jpg"
)

func NewAWSS3Config(credentials config.Credentials) AWSS3Config {
	return AWSS3Config{
		AccessKeyID:  credentials.AWSAccessKeyID,
		SecretKey:    credentials.AWSSecretKey,
		S3Region:     credentials.AWSS3ProfileRegion,
		S3BucketName: credentials.AWSS3ProfileBucket,
	}
}

func NewAWSS3Session(config AWSS3Config) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(config.S3Region),
		Credentials: credentials.NewStaticCredentials(
			config.AccessKeyID,
			config.SecretKey,
			"",
		),
	})
}

func NewS3StorageSession(credentials config.Credentials) (services.ImageStorageSession, error) {
	config := NewAWSS3Config(credentials)
	AWSSession, err := NewAWSS3Session(config)
	if err != nil {
		return nil, err
	}

	return S3StorageSession{
		awsS3Config: config,
		awsSession:  AWSSession,
	}, nil
}

type AWSS3Config struct {
	AccessKeyID  string
	SecretKey    string
	S3Region     string
	S3BucketName string
}

type S3StorageSession struct {
	awsS3Config AWSS3Config
	awsSession  *session.Session
}

func (object S3StorageSession) Upload(imageReader io.Reader, filename string) (string, error) {
	uploader := s3manager.NewUploader(object.awsSession)
	config := object.awsS3Config

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

func (object S3StorageSession) Download(imageIDParam string) (*services.ImageBufferedFile, error) {
	imageID := strings.ReplaceAll(imageIDParam, " ", "")
	if len(imageID) <= 0 {
		imageID = DEFAULT_IMAGE
	}

	downloader := s3manager.NewDownloader(object.awsSession)
	s3object := &s3.GetObjectInput{
		Bucket: aws.String(object.awsS3Config.S3BucketName),
		Key:    aws.String(imageID),
	}

	myFile := &aws.WriteAtBuffer{}
	bytesDownloaded, err := downloader.Download(myFile, s3object)
	if err != nil {
		return nil, response.MakeErrorJSON(http.StatusBadGateway, err.Error())
	}

	data := myFile.Bytes()

	return &services.ImageBufferedFile{Data: data, Size: bytesDownloaded}, nil
}

func (object S3StorageSession) Delete(objectID string) error {
	svc := s3.New(object.awsSession)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(object.awsS3Config.S3BucketName), Key: aws.String(objectID)})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(object.awsS3Config.S3BucketName),
		Key:    aws.String(objectID),
	})

	return err
}

func buildS3ImageURI(imageID string, bucketName string) string {
	url := url.URL{
		Scheme: S3_URI_SCHEME,
		Host:   bucketName,
		Path:   imageID,
	}

	return url.Path
}
