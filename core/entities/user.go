package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	User(email string) (User, error)
	Users() ([]User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(email string) error
	ExistUser(email string) bool
}

type User struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Email       string    `json:"email" form:"email" db:"email"`
	Nickname    string    `json:"nickname" form:"nickname" db:"nickname"`
	Password    string    `json:"password" form:"password" db:"password"`
	ImageID     string    `json:"image_id" db:"image_id"`
	CountryCode string    `json:"country_code" form:"country_code" db:"country_code"`
	Birthday    string    `json:"birthday" form:"birthday" db:"birthday"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (oldUser *User) Update(updatedUser User) {
	oldUser.Nickname = updatedUser.Nickname
	oldUser.Password = updatedUser.Password
	oldUser.ImageID = updatedUser.ImageID
	oldUser.CountryCode = updatedUser.CountryCode
	oldUser.Birthday = updatedUser.Birthday
}
