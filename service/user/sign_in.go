package user

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
	"website-api/common"
	"website-api/database/encrypt"
	authModel "website-api/model/auth"
	"website-api/model/role"
	userModel "website-api/model/user"
	"website-api/utils"
)

func (s *service) SignIn(reqBody *userModel.SignInRequest) (resData userModel.SignInResponse, statusCode int, err error) {
	if err := utils.CheckVerificationRateLimit("signin_" + reqBody.Email); err != nil {
		log.Printf("SECURITY: Sign in rate limit exceeded for email: %s", reqBody.Email)
		return resData, http.StatusTooManyRequests, err
	}

	userDB, err := s.userRepo.Take([]string{"id", "name", "user_name", "email", "encrypted_password", "phone_number", "is_verified", "role_id"}, &userModel.User{Email: reqBody.Email})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("SECURITY: Sign in attempt with non-existent email: %s", reqBody.Email)
			return resData, http.StatusUnauthorized, fmt.Errorf(common.EmailOrPasswordIsIncorrect)
		}
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	roleDB, err := s.roleRepo.Take([]string{"id", "name", "display_name", "description", "is_active"}, &role.Role{Id: userDB.RoleId})
	if err != nil {
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal mengambil data role: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDB.EncryptedPassword), []byte(reqBody.Password)); err != nil {
		log.Printf("SECURITY: Invalid password attempt for email: %s", reqBody.Email)
		return resData, http.StatusUnauthorized, fmt.Errorf(common.EmailOrPasswordIsIncorrect)
	}

	if !userDB.IsVerified {
		return resData, http.StatusForbidden, fmt.Errorf("email belum diverifikasi. Silakan cek email Anda untuk verifikasi")
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	claims := authModel.JWTClaims{
		UserID:   userDB.Id,
		Email:    userDB.Email,
		Username: userDB.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "website-api",
			Subject:   userDB.Id,
		},
	}

	jwtToken, err := encrypt.NewTokenWithClaims(claims)
	if err != nil {
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal membuat token: %w", err)
	}

	utils.ResetVerificationRateLimit("signin_" + reqBody.Email)

	resData = userModel.SignInResponse{
		Token: jwtToken,
		User: userModel.UserProfile{
			ID:          userDB.Id,
			Name:        userDB.Name,
			Username:    userDB.UserName,
			Email:       userDB.Email,
			PhoneNumber: userDB.PhoneNumber,
			IsVerified:  userDB.IsVerified,
			Role:        roleDB,
		},
		ExpiresAt: expiresAt.Unix(),
	}

	log.Printf("INFO: Successful sign in for email: %s", reqBody.Email)
	return resData, http.StatusOK, nil
}
