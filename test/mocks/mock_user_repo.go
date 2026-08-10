package mocks

import (
	userModel "website-api/model/user"

	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(data *userModel.User) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockUserRepo) Find(reqQuery *userModel.ListUserReqQuery) ([]userModel.ListUserResponse, int64, error) {
	args := m.Called(reqQuery)

	var users []userModel.ListUserResponse

	if args.Get(0) != nil {
		users = args.Get(0).([]userModel.ListUserResponse)
	}

	return users, args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepo) Take(selects []string, conditions *userModel.User) (userModel.User, error) {

	args := m.Called(selects, conditions)

	if args.Get(0) == nil {
		return userModel.User{}, args.Error(1)
	}

	return args.Get(0).(userModel.User), args.Error(1)
}

func (m *MockUserRepo) Update(id *string, data *map[string]interface{}) error {
	args := m.Called(id, data)
	return args.Error(0)
}

func (m *MockUserRepo) Seed() error {
	args := m.Called()
	return args.Error(0)
}
