package entities

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID `json:"id" form:"id" db:"id"`
	Email       string    `json:"email" form:"email" db:"email"`
	Nickname    string    `json:"nickname" form:"nickname" db:"nickname"`
	Password    string    `json:"password" form:"password" db:"password"`
	ImageURL    string    `json:"image_url" form:"image_url" db:"image_url"`
	CountryCode string    `json:"country_code" form:"country_code" db:"country_code"`
	Birthday    string    `json:"birthday" form:"birthday" db:"birthday"`
}

type UserRepository interface {
	User(email string) (User, error)
	Users() ([]User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(email string) error
}

type Repository interface {
	UserRepository
}
