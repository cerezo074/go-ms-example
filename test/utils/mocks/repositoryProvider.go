package mocks

import . "user/core/entities"

type FakeRepo struct {
	AllUsers    func() ([]User, error)
	UserByEmail func(string) (User, error)
	Save        func(*User) error
	Update      func(*User) error
	Delete      func(string) error
	Exist       func(string) bool
}

func (object FakeRepo) User(email string) (User, error) {
	if object.UserByEmail == nil {
		return User{}, nil
	}

	return object.UserByEmail(email)
}

func (object FakeRepo) Users() ([]User, error) {
	if object.AllUsers == nil {
		return []User{}, nil
	}

	return object.AllUsers()
}

func (object FakeRepo) CreateUser(user *User) error {
	if object.Save == nil {
		return nil
	}

	return object.Save(user)
}

func (object FakeRepo) UpdateUser(user *User) error {
	if object.Update == nil {
		return nil
	}

	return object.Update(user)
}

func (object FakeRepo) DeleteUser(email string) error {
	if object.Delete == nil {
		return nil
	}

	return object.Delete(email)
}

func (object FakeRepo) ExistUser(email string) bool {
	if object.Exist == nil {
		return false
	}

	return object.Exist(email)
}
