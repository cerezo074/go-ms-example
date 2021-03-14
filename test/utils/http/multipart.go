package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"reflect"
	"strings"
)

type UserForm struct {
	Email       string  `form:"email"`
	Nickname    string  `form:"nickname"`
	Password    string  `form:"password"`
	ImagePath   *string `form:"image_data" type:"file"`
	CountryCode string  `form:"country_code"`
	Birthday    string  `form:"birthday"`
}

func (object UserForm) imageName() string {
	_, filename := path.Split(*object.ImagePath)
	return filename
}

func FilesMatch(rawFile io.Reader, rigthFilePath string) (bool, error) {
	rawImage, err := os.Open(rigthFilePath)
	if err != nil {
		return false, err
	}

	fileInfo, err := rawImage.Stat()
	if err != nil {
		return false, err
	}

	bytes := bytes.NewBuffer([]byte{})
	bytesWrited, err := io.Copy(bytes, rawFile)
	if err != nil {
		return false, err
	}

	if fileInfo.Size() != bytesWrited {
		return false, errors.New("Files dont contain same size")
	}

	return true, nil
}

func AddMultipartFile(key string, filepath string, writer *multipart.Writer) error {
	_, filename := path.Split(filepath)
	if len(filename) == 0 || strings.Contains(filename, " ") {
		return errors.New(fmt.Sprintf("Invalid filename %v", filename))
	}

	part, err := writer.CreateFormFile(key, filename)
	if err != nil {
		return err
	}

	sample, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer sample.Close()
	_, err = io.Copy(part, sample)
	if err != nil {
		return err
	}

	err = writer.Close()
	return nil
}

func AddMultipartField(key string, value string, writer *multipart.Writer) error {
	field, err := writer.CreateFormField(key)
	if err != nil {
		return err
	}

	_, err = field.Write([]byte(value))
	return nil
}

func MultipartFormBody(form interface{}) (body *bytes.Buffer, contentType string, bodyErr error) {
	defer func() {
		err := recover()
		panicError, ok := err.(error)
		if ok {
			body = nil
			bodyErr = fmt.Errorf("Invalid type form %t, needs to be a pointer or a interface type. Error, %v", form, panicError)
			contentType = ""
		}
	}()

	body = new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	value := reflect.ValueOf(form).Elem()

	for i := 0; i < value.NumField(); i++ {
		typeField := value.Type().Field(i)
		valueField := value.Field(i)
		tag := typeField.Tag

		formKey := tag.Get("form")
		if formKey == "" {
			continue
		}

		formValue := valueField.Interface()
		var stringValue string
		if stringValueReference, okValue := formValue.(string); okValue {
			stringValue = stringValueReference
		} else if stringPointerReference, okPointer := formValue.(*string); okPointer && stringPointerReference != nil {
			stringValue = *stringPointerReference
		} else {
			continue
		}

		if tag.Get("type") == "file" {
			err := AddMultipartFile(formKey, stringValue, writer)
			if err != nil {
				return nil, "", err
			}
		} else {
			err := AddMultipartField(formKey, stringValue, writer)
			if err != nil {
				return nil, "", err
			}
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}
