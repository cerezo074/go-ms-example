package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id" form:"id" db:"id"`
	Email       string    `json:"email" form:"email" db:"email"`
	Nickname    string    `json:"nickname" form:"nickname" db:"nickname"`
	Password    string    `json:"password" form:"password" db:"password"`
	ImageURL    string    `json:"image_url" form:"image_url" db:"image_url"`
	CountryCode string    `json:"country_code" form:"country_code" db:"country_code"`
	Birthday    string    `json:"birthday" form:"birthday" db:"birthday"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (oldUser *User) Update(updatedUser User) {
	oldUser.Nickname = updatedUser.Nickname
	oldUser.Password = updatedUser.Password
	oldUser.ImageURL = updatedUser.ImageURL
	oldUser.CountryCode = updatedUser.CountryCode
	oldUser.Birthday = updatedUser.Birthday
}

func (user User) IsValid() error {
	errorMessage := ""
	if user.ID.String() == "" {
		errorMessage += "Invalid id, "
	}

	if user.Email == "" {
		errorMessage += "Invalid email, "
	}

	if user.Nickname == "" {
		errorMessage += "Invalid nickname, "
	}

	if user.Password == "" {
		errorMessage += "Invalid password, "
	}

	if user.ImageURL == "" {
		errorMessage += "Invalid image url, "
	}

	if user.CountryCode == "" {
		errorMessage += "Invalid country code, "
	}

	if user.Birthday == "" {
		errorMessage += "Invalid birthday"
	}

	if errorMessage == "" {
		return nil
	}

	return errors.New("There are some fields invalids: " + errorMessage)
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
