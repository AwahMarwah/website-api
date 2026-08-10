package mocks

import (
	authModel "website-api/model/auth"

	"github.com/stretchr/testify/mock"
)

type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) Create(auth *authModel.PasswordResetToken) error {
	args := m.Called(auth)

	return args.Error(0)
}

func (m *MockAuthRepo) Take(selectParams []string, condition *authModel.PasswordResetToken) (authModel.PasswordResetToken, error) {
	args := m.Called(selectParams, condition)

	var result authModel.PasswordResetToken

	if args.Get(0) != nil {
		result = args.Get(0).(authModel.PasswordResetToken)
	}
	return result, args.Error(1)
}

func (m *MockAuthRepo) Update(id *string, values *map[string]any) error {
	args := m.Called(id, values)
	return args.Error(0)
}

func (m *MockAuthRepo) Delete(condition map[string]interface{}) error {
	args := m.Called(condition)
	return args.Error(0)
}
