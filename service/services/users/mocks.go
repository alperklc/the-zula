package usersService

import (
	"github.com/stretchr/testify/mock"
)

type MockedUser struct {
	mock.Mock
}

func (m *MockedUser) GetUser(id string) (User, error) {
	args := m.Called(id)
	return args.Get(0).(User), args.Error(1)
}

func (m *MockedUser) UpdateUser(id, firstName, lastname, displayName string, language, theme *string) error {
	args := m.Called(id, firstName, lastname, displayName, language, theme)
	return args.Error(0)
}
