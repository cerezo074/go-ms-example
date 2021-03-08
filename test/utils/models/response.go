package models

type AllUsersResponse struct {
	Data []User `json:"data"`
}

type FindUserResponse struct {
	Data User `json:"data"`
}

type FailResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data string `json:"data"`
}
