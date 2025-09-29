package user

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"net/http"
	"os"
	"testing"
	userModel "website-api/model/user"
)

// Setup test environment
func setupTestEnv() {
	os.Setenv("APP_BASE_URL", "http://localhost:8080")
	os.Setenv("JWT_SECRET", "test_secret")
	os.Setenv("SMTP_HOST", "smtp.test.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USERNAME", "test@example.com")
	os.Setenv("SMTP_PASSWORD", "testpassword")
	os.Setenv("TEMPTLATE_PATH", "../../templates")

}

// Cleanup test environment
func cleanupTestEnv() {
	os.Unsetenv("APP_BASE_URL")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("SMTP_USERNAME")
	os.Unsetenv("SMTP_PASSWORD")
}

// MockUserRepo tetap sama seperti sebelumnya
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Take(columns []string, conditions *userModel.User) (userModel.User, error) {
	args := m.Called(columns, conditions)
	return args.Get(0).(userModel.User), args.Error(1)
}

func (m *MockUserRepo) Create(user *userModel.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepo) Find(reqQuery *userModel.ListUserReqQuery) (users []userModel.ListUserResponse, count int64, err error) {
	args := m.Called(reqQuery)
	return args.Get(0).([]userModel.ListUserResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepo) Update(id *string, values *map[string]any) error {
	args := m.Called(id, values)
	return args.Error(0)
}

// Test helper untuk membuat request yang valid
func createValidRequest() *userModel.RegisterRequest {
	return &userModel.RegisterRequest{
		Name:        "John Doe",
		Username:    "johndoe",
		Email:       "johndoe@example.com",
		Password:    "password123",
		PhoneNumber: "1234567890",
	}
}

// Test sign_up_success dengan environment setup
func TestSignUpSuccess(t *testing.T) {
	// Setup environment
	setupTestEnv()
	defer cleanupTestEnv()

	// Setup
	mockRepo := new(MockUserRepo)
	service := &service{userRepo: mockRepo}
	reqBody := createValidRequest()

	// Mock: email not found in database
	mockRepo.On("Take", []string{"email"}, &userModel.User{Email: reqBody.Email}).
		Return(userModel.User{}, gorm.ErrRecordNotFound)

	// Mock: create successfully user
	mockRepo.On("Create", mock.AnythingOfType("*user.User")).Return(nil)

	// Execute
	statusCode, err := service.SignUp(reqBody)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, statusCode)
	mockRepo.AssertExpectations(t)
}

// Test lainnya tetap sama, tapi tambahkan setupTestEnv() jika diperlukan
func TestSignUpEmailAlreadyExists(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepo)
	service := NewService(mockRepo)
	reqBody := createValidRequest()

	// Mock: email already exist
	existingUser := userModel.User{Email: reqBody.Email}
	mockRepo.On("Take", []string{"email"}, &userModel.User{Email: reqBody.Email}).
		Return(existingUser, nil)

	// Execute
	statusCode, err := service.SignUp(reqBody)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Contains(t, err.Error(), "sudah terdaftar")
	mockRepo.AssertExpectations(t)
}

func TestSignUpDatabaseError(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepo)
	service := NewService(mockRepo)
	reqBody := createValidRequest()

	// Mock: error when checking email in database
	mockRepo.On("Take", []string{"email"}, &userModel.User{Email: reqBody.Email}).
		Return(userModel.User{}, gorm.ErrInvalidDB)

	// Execute
	statusCode, err := service.SignUp(reqBody)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Contains(t, err.Error(), "failed to check existing user")
	mockRepo.AssertExpectations(t)
}

func TestSignUpCreateUserError(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepo)
	service := NewService(mockRepo)
	reqBody := createValidRequest()

	// Mock: email not found, but error when creating user
	mockRepo.On("Take", []string{"email"}, &userModel.User{Email: reqBody.Email}).
		Return(userModel.User{}, gorm.ErrRecordNotFound)

	mockRepo.On("Create", mock.AnythingOfType("*user.User")).Return(gorm.ErrInvalidDB)

	// Execute
	statusCode, err := service.SignUp(reqBody)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Contains(t, err.Error(), "failed to create user")
	mockRepo.AssertExpectations(t)
}
