package request

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserResponse struct {
	Data []User `json:"data"`
}

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Nickname    string `json:"nickname"`
	Password    string `json:"password"`
	ImageID     string `json:"image_id"`
	CountryCode string `json:"country_code"`
	Birthday    string `json:"birthday"`
	CreatedAt   string `json:"CreatedAt"`
	UpdatedAt   string `json:"UpdatedAt"`
}

type FakeServer struct {
	FiberApp *fiber.App
}

func (object FakeServer) GetJSONObject(response *http.Response) (UserResponse, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return UserResponse{}, err
	}

	var userResponse UserResponse
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return UserResponse{}, err
	}

	return userResponse, nil
}

func (object FakeServer) Execute(method string, path string, body io.Reader, dataKey string) (*http.Response, []User, error) {
	request, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, nil, err
	}

	response, err := object.FiberApp.Test(request, -1)
	if err != nil {
		return nil, nil, err
	}

	userResponse, err := object.GetJSONObject(response)
	if err != nil {
		return nil, nil, err
	}

	if len(dataKey) > 0 && dataKey != " " {
		return response, userResponse.Data, nil
	}

	return nil, nil, fmt.Errorf("Invalid JSONSlice in object, %v", userResponse)
}
