package models

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
