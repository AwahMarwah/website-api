package user

import (
	"errors"
	"fmt"
	"log"
	"time"
	"website-api/common"
	userModel "website-api/model/user"
	"website-api/utils"

	"gorm.io/gorm"
)

func (s *service) VerifyEmail(reqBody userModel.VerifyEmailRequest) (err error) {
	if !utils.ValidateTokenFormat(reqBody.Token) {
		log.Printf("SECURITY: Invalid token format attempted: %s", tokenPreview(reqBody.Token))
		return fmt.Errorf(common.InvalidVerificationToken)
	}

	if err := utils.CheckVerificationRateLimit(reqBody.Token); err != nil {
		log.Printf("SECURITY: Rate limit exceeded for token: %s", tokenPreview(reqBody.Token))
		return err
	}

	user, err := s.userRepo.Take([]string{"id", "email", "verification_token", "verification_token_expired_at", "is_verified"}, &userModel.User{VerificationToken: reqBody.Token})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("SECURITY: Invalid verification token attempted: %s", tokenPreview(reqBody.Token))
			return fmt.Errorf(common.InvalidVerificationToken)
		}
		log.Printf("SECURITY: Database error during verification: %v", err)
		return fmt.Errorf(common.VerificationFailed)
	}

	if user.Id == "" || user.VerificationToken == "" {
		log.Printf("SECURITY: Empty user data for token: %s", tokenPreview(reqBody.Token))
		return fmt.Errorf(common.InvalidVerificationToken)
	}

	if time.Now().After(user.VerificationTokenExpiredAt) {
		log.Printf("SECURITY: Expired token attempted: %s for user: %s", tokenPreview(reqBody.Token), user.Email)
		values := map[string]any{
			"verification_token":            "",
			"verification_token_expired_at": nil,
			"updated_at":                    time.Now(),
			"updated_by":                    "system-cleanup",
		}
		if err = s.userRepo.Update(&user.Id, &values); err != nil {
			return fmt.Errorf("failed to update user verification status: %w", err)
		}
		return fmt.Errorf(common.VerificationTokenExpired)
	}

	if user.IsVerified {
		log.Printf("SECURITY: Already verified email attempted: %s", user.Email)
		return fmt.Errorf(common.EmailAlreadyVerified)
	}

	values := map[string]any{
		"is_verified":                   true,
		"verification_token":            "",
		"verification_token_expired_at": nil,
		"updated_at":                    time.Now(),
		"updated_by":                    "email-verification",
	}

	if err := s.userRepo.Update(&user.Id, &values); err != nil {
		return fmt.Errorf("failed to update user verification status: %w", err)
	}

	utils.ResetVerificationRateLimit(reqBody.Token)

	return nil
}

// tokenPreview mengembalikan 8 karakter awal token untuk keperluan log
// keamanan tanpa membocorkan token penuh dan tanpa panic pada token pendek.
func tokenPreview(token string) string {
	if len(token) > 8 {
		return token[:8] + "..."
	}
	return token
}
