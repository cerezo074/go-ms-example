package postgres

import (
	"errors"
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
		return entities.User{}, fmt.Errorf("fail to retrieve one user, %w", err)
	}

	return user, nil
}

func (store *UserStore) Users() ([]entities.User, error) {
	var users []entities.User

	if err := store.Select(&users, "SELECT * FROM users"); err != nil {
		return []entities.User{}, fmt.Errorf("fail to retrieve all users, %w", err)
	}

	return users, nil
}

func (store *UserStore) CreateUser(newUser *entities.User) (execError error) {
	defer func() {
		if err := recover(); err != nil {
			execError = store.processRecover(err)
		}
	}()

	const query = "INSERT INTO users (email, nickname, password, image_id, country_code, birthday) VALUES ($1, $2, $3, $4, $5, $6)"

	result := store.MustExec(query,
		newUser.Email,
		newUser.Nickname,
		newUser.Password,
		newUser.ImageID,
		newUser.CountryCode,
		newUser.Birthday)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fail to create a new user, %w", err)
	}

	if rowsAffected == 1 {
		return nil
	}

	return errors.New("Invalid rows affected")
}

func (store *UserStore) UpdateUser(oldUser *entities.User) (execError error) {
	defer func() {
		if err := recover(); err != nil {
			execError = store.processRecover(err)
		}
	}()

	const query = "UPDATE users SET nickname = $1, password = $2, image_uri = $3, country_code = $4, birthday = $5 WHERE email = $6"

	result := store.MustExec(query,
		oldUser.Nickname,
		oldUser.Password,
		oldUser.ImageID,
		oldUser.CountryCode,
		oldUser.Birthday,
		oldUser.Email)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fail to update a user, %w", err)
	}

	if rowsAffected == 1 {
		return nil
	}

	return errors.New("Invalid rows affected")
}

func (store *UserStore) DeleteUser(email string) (execError error) {
	defer func() {
		if err := recover(); err != nil {
			execError = store.processRecover(err)
		}
	}()

	result := store.MustExec("DELETE FROM users WHERE email = $1", email)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("fail deleting a user, %w", err)
	}

	if rowsAffected == 1 {
		return nil
	}

	return errors.New("Invalid rows affected")
}

func (store *UserStore) ExistUser(email string) bool {
	var rows int

	if err := store.Get(&rows, "SELECT count(*) FROM users WHERE email = $1", email); err != nil {
		return false
	}

	return rows == 1
}

func (store *UserStore) processRecover(value interface{}) error {
	switch err := value.(type) {
	case string:
		return errors.New(err)
	case error:
		return err
	default:
		return errors.New("Unknown panic")
	}
}
