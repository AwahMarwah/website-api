package auth

import (
	"errors"
	"fmt"
	"net/http"
	"website-api/common"
	authModel "website-api/model/auth"
	userModel "website-api/model/user"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *service) ChangePassword(userID string, req *authModel.ChangePasswordRequest) (int, error) {
	user, err := s.userRepo.Take(
		[]string{"id", "encrypted_password"},
		&userModel.User{Id: userID},
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, fmt.Errorf("pengguna tidak ditemukan")
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mengambil pengguna: %w", err)
	}

	// Verifikasi password lama
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.EncryptedPassword),
		[]byte(req.CurrentPassword),
	); err != nil {
		return http.StatusBadRequest, fmt.Errorf(common.EmailOrPasswordIsIncorrect)
	}

	// Hash password baru
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mengenkripsi password: %w", err)
	}

	// Update di database
	if err := s.userRepo.Update(&user.Id, &map[string]any{
		"encrypted_password": string(newHash),
	}); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memperbarui password: %w", err)
	}

	return http.StatusOK, nil
}