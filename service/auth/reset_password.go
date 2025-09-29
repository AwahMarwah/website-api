package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"time"
	authModel "website-api/model/auth"
)

func (s *service) ResetPassword(req *authModel.ResetPasswordRequest) (statusCode int, message string, err error) {
	// Hash token untuk mencari di database
	tokenHash := sha256.Sum256([]byte(req.Token))
	tokenHashStr := hex.EncodeToString(tokenHash[:])

	// Cari token di database
	resetToken, err := s.authRepo.Take([]string{"id", "user_id", "expires_at", "is_used"}, &authModel.PasswordResetToken{TokenHash: tokenHashStr})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusBadRequest, "", fmt.Errorf("token reset password tidak valid")
		}
		return http.StatusInternalServerError, "", fmt.Errorf("gagal memverifikasi token: %w", err)
	}

	// Cek apakah token sudah digunakan
	if resetToken.IsUsed {
		return http.StatusBadRequest, "", fmt.Errorf("token reset password sudah digunakan")
	}

	// Cek apakah token sudah expired
	if time.Now().After(resetToken.ExpiresAt) {
		return http.StatusBadRequest, "", fmt.Errorf("token reset password sudah kedaluwarsa")
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("gagal mengenkripsi password: %w", err)
	}

	userValues := map[string]any{
		"encrypted_password": string(hashedPassword),
		"updated_at":         time.Now(),
		"updated_by":         "reset-password",
	}

	if err := s.userRepo.Update(&resetToken.UserID, &userValues); err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("gagal mengupdate password: %w", err)
	}

	tokenValues := map[string]any{
		"is_used": true,
	}

	if err := s.authRepo.Update(&resetToken.Id, &tokenValues); err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("gagal mengupdate token: %w", err)
	}

	return http.StatusOK, "Password berhasil direset", nil
}
