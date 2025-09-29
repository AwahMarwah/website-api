package user

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
	"website-api/common"
	lib "website-api/library/helper/email"
	roleModel "website-api/model/role"
	userModel "website-api/model/user"
	"website-api/utils"
)

func (s *service) SignUp(reqBody *userModel.RegisterRequest) (statusCode int, err error) {
	if err := utils.CheckVerificationRateLimit("signup_" + reqBody.Email); err != nil {
		log.Printf("SECURITY: Sign up rate limit exceeded for email: %s", reqBody.Email)
		return http.StatusTooManyRequests, err
	}
	userDB, err := s.userRepo.Take([]string{"email"}, &userModel.User{Email: reqBody.Email})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusInternalServerError, fmt.Errorf("failed to check existing user: %w", err)
	}
	if userDB.Email != "" {
		return http.StatusBadRequest, fmt.Errorf(common.EmailAlreadyExists)
	}

	if reqBody.RoleName == "" {
		reqBody.RoleName = "consumer"
	}

	roleData, err := s.roleRepo.Take([]string{"id", "name", "display_name"}, &roleModel.Role{Name: reqBody.RoleName})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusInternalServerError, fmt.Errorf("failed to check role: %w", err)
	}
	if roleData.Id == "" {
		return http.StatusBadRequest, fmt.Errorf("role %s tidak ditemukan", reqBody.RoleName)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error saat hashing password: %v\n", err)
		return http.StatusInternalServerError, fmt.Errorf("failed to hash password: %w", err)
	}

	loginToken := uuid.NewString()
	verificationToken, err := utils.GenerateVerificationToken()
	expiredAt := time.Now().Add(24 * time.Hour)

	user := userModel.User{
		Name:                       reqBody.Name,
		UserName:                   reqBody.Username,
		Email:                      reqBody.Email,
		EncryptedPassword:          string(hashedPassword),
		Token:                      loginToken,
		PhoneNumber:                reqBody.PhoneNumber,
		VerificationToken:          verificationToken,
		VerificationTokenExpiredAt: expiredAt,
		IsVerified:                 false,
		CreatedAt:                  time.Now(),
		CreatedBy:                  "system",
		RoleId:                     roleData.Id,
	}

	if err := s.userRepo.Create(&user); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create user: %w", err)
	}

	if err := lib.SendVerificationEmail(&user, verificationToken); err != nil {
		log.Printf("SECURITY: Failed to send verification email to %s: %v", user.Email, err)
	}
	return http.StatusCreated, nil
}
