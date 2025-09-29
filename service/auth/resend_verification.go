package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	lib "website-api/library/helper/email"
	authModel "website-api/model/auth"
	userModel "website-api/model/user"
	"website-api/utils"

	"gorm.io/gorm"
)

func (s *service) ResendVerification(req *authModel.ResendVerificationRequest) (statusCode int, err error) {
	user, err := s.userRepo.Take([]string{"id", "name", "email", "is_verified"}, &userModel.User{Email: req.Email})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusBadRequest, fmt.Errorf("email tidak terdaftar")
		}
		return http.StatusInternalServerError, fmt.Errorf("failed to find user: %w", err)
	}

	if user.IsVerified {
		return http.StatusBadRequest, fmt.Errorf("email sudah diverifikasi")
	}

	// Generate token baru
	verificationToken, err := utils.GenerateVerificationToken()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to generate verification token: %w", err)
	}

	expiredAt := time.Now().Add(24 * time.Hour)

	values := map[string]any{
		"verification_token":            verificationToken,
		"verification_token_expired_at": expiredAt,
		"updated_at":                    time.Now(),
		"updated_by":                    "resend-verification",
	}

	if err := s.userRepo.Update(&user.Id, &values); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to update verification token: %w", err)
	}

	if err := lib.SendVerificationEmail(&user, verificationToken); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to send verification email: %w", err)
	}

	return http.StatusOK, nil
}
