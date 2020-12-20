package postgres

import (
	"fmt"
	"user/core/entities"

	"github.com/jmoiron/sqlx"
)

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{
		DB: db,
	}
}

type UserStore struct {
	*sqlx.DB
}

func (store *UserStore) User(email string) (entities.User, error) {
	var user entities.User

	if err := store.Get(&user, "SELECT * FROM users WHERE email = $1", email); err != nil {
		return entities.User{}, fmt.Errorf("error retrieving one user: %w", err)
	}

	return user, nil
}

func (store *UserStore) Users() ([]entities.User, error) {
	var users []entities.User

	if err := store.Select(&users, "SELECT * FROM users"); err != nil {
		return []entities.User{}, fmt.Errorf("error retrieving all users: %w", err)
	}

	return users, nil
}

func (store *UserStore) CreateUser(newUser *entities.User) error {
	if err := store.Get(newUser, "INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6)",
		newUser.Email,
		newUser.Nickname,
		newUser.Password,
		newUser.ImageURL,
		newUser.CountryCode,
		newUser.Birthday); err != nil {
		return fmt.Errorf("error creating a new user: %w", err)
	}

	return nil
}

func (store *UserStore) UpdateUser(oldUser *entities.User) error {
	const query = "UPDATE users SET nickname = $1, password = $2, image_url = $3, country_code = $4, birthday = $5 WHERE email = $6"

	if err := store.Get(oldUser, query,
		oldUser.Nickname,
		oldUser.Password,
		oldUser.ImageURL,
		oldUser.CountryCode,
		oldUser.Birthday,
		oldUser.Email); err != nil {
		return fmt.Errorf("error updating a user: %w", err)
	}

	return nil
}

func (store *UserStore) DeleteUser(email string) error {
	if _, err := store.Exec("DELETE FROM users WHERE email = $1", email); err != nil {
		return fmt.Errorf("error deleting a user: %w", err)
	}

	return nil
}
