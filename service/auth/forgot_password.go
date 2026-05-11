package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"website-api/common"
	authModel "website-api/model/auth"
	userModel "website-api/model/user"
	"website-api/task"
	"website-api/worker"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

func (s *service) ForgotPassword(req *authModel.ForgotPasswordRequest) (statusCode int, message string, err error) {
	user, err := s.userRepo.Take([]string{"id", "name", "email"}, &userModel.User{Email: req.Email})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusOK, common.PasswordResetRequestSent, nil
		}
		return http.StatusInternalServerError, "", fmt.Errorf("gagl mencari pengguna: %w", err)
	}

	// Hapus token reset password yang lama jika ada
	err = s.authRepo.Update(&user.Id, &map[string]interface{}{
		"user_id": user.Id,
		"is_used": false,
	})

	// Generate token
	token := uuid.NewString()
	tokenHash := sha256.Sum256([]byte(token))
	tokenHashStr := hex.EncodeToString(tokenHash[:])

	// Assigne to model
	resetRecord := &authModel.PasswordResetToken{
		Id:        uuid.NewString(),
		UserID:    user.Id,
		TokenHash: tokenHashStr,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		IsUsed:    false,
	}

	// Save to database
	if err := s.authRepo.Create(resetRecord); err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("gagl menyimpan token reset password: %w", err)
	}

	// send email reset password
	// TO DO: replace with actual email service
	//if err := lib.SendResetPasswordByEmail(&user, token); err != nil {
	//	log.Printf("Failed to send reset password email to %s: %v", user.Email, err)
	//}

	redisClient := worker.NewRedisClient()
	defer redisClient.Close()

	emailTask, err := task.NewResetPasswordTask(user.Name, user.Email, token)
	if err != nil {
		log.Printf("failed create task: %v", err)
	} else {
		_, err = redisClient.Enqueue(
			emailTask,
			asynq.MaxRetry(3),
			asynq.Queue("critical"),
		)
		if err != nil {
			log.Printf("failed enqueue email: %v", err)
		}
	}
	return http.StatusOK, common.PasswordResetRequestSent, nil
}
