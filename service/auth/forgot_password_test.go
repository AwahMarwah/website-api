package auth

import (
	"net/http"
	"testing"
	"website-api/common"
	authModel "website-api/model/auth"
	userModel "website-api/model/user"
	"website-api/test/mocks"

	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestForgotPassword(t *testing.T) {
	userRepo := new(mocks.MockUserRepo)
	authRepo := new(mocks.MockAuthRepo)
	queueClient := new(mocks.MockQueueClient)

	userRepo.On("Take", mock.Anything, mock.Anything).Return(userModel.User{
		Id:    "user-1",
		Name:  "Asep",
		Email: "a@a.com",
	}, nil)

	authRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	authRepo.On("Create", mock.Anything).Return(nil)

	queueClient.On("Enqueue", mock.Anything).Return(&asynq.TaskInfo{}, nil)

	svc := NewService(userRepo, authRepo, queueClient)
	status, message, err := svc.ForgotPassword(&authModel.ForgotPasswordRequest{Email: "a@a.com"})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, common.PasswordResetRequestSent, message)

	userRepo.AssertExpectations(t)
	authRepo.AssertExpectations(t)
	queueClient.AssertExpectations(t)
}

func TestForgotPasswordUserNotFound(t *testing.T) {
	userRepo := new(mocks.MockUserRepo)
	authRepo := new(mocks.MockAuthRepo)
	queueClient := new(mocks.MockQueueClient)

	userRepo.On("Take", mock.Anything, mock.Anything).Return(userModel.User{}, gorm.ErrRecordNotFound)

	svc := NewService(userRepo, authRepo, queueClient)
	status, message, err := svc.ForgotPassword(&authModel.ForgotPasswordRequest{Email: "usernotfound@gmail.com"})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, common.PasswordResetRequestSent, message)

	userRepo.AssertExpectations(t)
}

func TestForgotPasswordUserRepoError(t *testing.T) {
	userRepo := new(mocks.MockUserRepo)
	authRepo := new(mocks.MockAuthRepo)
	queueClient := new(mocks.MockQueueClient)

	userRepo.On("Take", mock.Anything, mock.Anything).Return(userModel.User{}, errors.New("database error"))

	svc := NewService(userRepo, authRepo, queueClient)
	status, _, err := svc.ForgotPassword(&authModel.ForgotPasswordRequest{Email: "a@a.com"})
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestForgotPasswordCreateTokenError(t *testing.T) {
	userRepo := new(mocks.MockUserRepo)
	authRepo := new(mocks.MockAuthRepo)
	queueClient := new(mocks.MockQueueClient)

	userRepo.On("Take", mock.Anything, mock.Anything).
		Return(userModel.User{
			Id:    "user-1",
			Name:  "Asep",
			Email: "z@z.com",
		}, nil)

	authRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	authRepo.On("Create", mock.Anything).Return(errors.New("insert error"))
	t.Log("CREATE MOCK REGISTERED")
	svc := NewService(userRepo, authRepo, queueClient)

	status, _, err := svc.ForgotPassword(&authModel.ForgotPasswordRequest{Email: "z@z.com"})
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestForgotPasswordQueueError(t *testing.T) {
	userRepo := new(mocks.MockUserRepo)
	authRepo := new(mocks.MockAuthRepo)
	queueClient := new(mocks.MockQueueClient)

	userRepo.On("Take", mock.Anything, mock.Anything).
		Return(userModel.User{
			Id:    "user-1",
			Name:  "Asep",
			Email: "z@z.com",
		}, nil)

	authRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	authRepo.On("Create", mock.Anything).Return(nil)
	queueClient.On("Enqueue", mock.Anything).Return(nil, errors.New("redis error"))

	svc := NewService(userRepo, authRepo, queueClient)
	status, message, err := svc.ForgotPassword(&authModel.ForgotPasswordRequest{Email: "z@z.com"})
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Equal(t, common.PasswordResetRequestSent, message)
}
