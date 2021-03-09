package request

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	. "user/test/utils/models"

	"github.com/gofiber/fiber/v2"
)

type ResponseUnmarshaller func(body []byte) (interface{}, error)

type FakeServer struct {
	FiberApp     *fiber.App
	Unmarshaller ResponseUnmarshaller
}

func NewAllUsersUnmarshaller(body []byte) (interface{}, error) {
	var response AllUsersResponse
	err := json.Unmarshal(body, &response)
	return response.Data, err
}

func NewFindUserUnmarshaller(body []byte) (interface{}, error) {
	var response FindUserResponse
	err := json.Unmarshal(body, &response)
	return response.Data, err
}

func NewFailUnmarshaller(body []byte) (interface{}, error) {
	var response FailResponse
	err := json.Unmarshal(body, &response)
	return response, err
}

func NewSuccessUnmarshaller(body []byte) (interface{}, error) {
	var response SuccessResponse
	err := json.Unmarshal(body, &response)
	return response, err
}

func (object FakeServer) GetJSONObject(response *http.Response) (interface{}, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if object.Unmarshaller == nil {
		return nil, errors.New("Invalid response marshaller")
	}

	data, err := object.Unmarshaller(body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (object FakeServer) Execute(method string, path string, body io.Reader) (*http.Response, interface{}, error) {
	request, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, nil, err
	}
	response, err := object.FiberApp.Test(request, -1)
	if err != nil {
		return nil, nil, err
	}

	responseData, err := object.GetJSONObject(response)
	if err != nil {
		return nil, nil, err
	}

	return response, responseData, nil
}
