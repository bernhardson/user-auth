package repo

import (
	"errors"
	"time"

	"github.com/bernhardson/stub/internal/models"
)

var mockUser = models.User{
	ID:       1,
	Username: "John",
	Email:    "john.doe@gmail.com",
	Password: "jd12345678",
	Created:  time.Now(),
}

type MockUserRepo struct {
}

func (m *MockUserRepo) Insert(name, email, password string) (int, error) {
	return 1, nil
}

func (m *MockUserRepo) Get(id int64) (*models.User, error) {
	switch id {
	case 1:
		return &mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *MockUserRepo) GetAll() (*[]models.User, error) {

	var mockUsers = make([]models.User, 1)
	mockUsers = append(mockUsers, mockUser)

	return &mockUsers, nil
}

func (m *MockUserRepo) DeleteUser(id int) (int, error) {
	if id != 1 {
		return id, errors.New("user not found")
	}
	return 1, nil
}

func (m *MockUserRepo) Authenticate(email, password string) (int, error) {
	if email == mockUser.Email && password == mockUser.Password {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}

func (m *MockUserRepo) Exists(id int) (bool, error) {
	if id == 1 {
		return true, nil
	}
	return false, nil
}

func (m *MockUserRepo) ClearTable(table string) error {
	return nil
}
