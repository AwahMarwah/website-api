package email

import (
	"errors"
	"fmt"
	"log"
	"os"
	userModel "website-api/model/user"
	"website-api/utils/email"
	"website-api/utils/template"
)

// Helper function untuk mengirim email verifikasi
func SendVerificationEmail(user *userModel.User, verificationToken string) error {
	subject := "Verifikasi Email Akun Anda - Website Simple Ecommerce"
	appBaseURL := os.Getenv("APP_BASE_URL")
	if appBaseURL == "" {
		return fmt.Errorf("APP_BASE_URL environment variable is not set")
	}

	verifyLink := fmt.Sprintf("%s/user/verify-email?token=%s", appBaseURL, verificationToken)

	body, err := template.RenderVerificationEmail(user.Name, verifyLink)
	if err != nil {
		log.Printf("Gagal render templates email: %v", err)
		return fmt.Errorf("failed to render email templates: %w", err)
	}

	emailSender := email.NewSMTPFromEnv()
	go func() {
		if err := emailSender.SendEmail(user.Email, subject, body); err != nil {
			log.Printf("Gagal mengirim email verifikasi ke %s: %v", user.Email, err)
		} else {
			log.Printf("Email verifikasi berhasil dikirim ke %s", user.Email)
		}
	}()

	return nil
}

func SendResetPasswordByEmail(user *userModel.User, token string) error {
	subject := "Permintaan Reset Password - Website Simple Ecommerce"
	appBaseURL := os.Getenv("APP_BASE_URL")
	if appBaseURL == "" {
		return errors.New("APP_BASE_URL environment variable is not set")
	}

	resetLink := fmt.Sprintf("%s/reset-password/%s", appBaseURL, token)

	body, err := template.RenderResetPasswordEmail(user.Name, resetLink)
	if err != nil {
		return fmt.Errorf("gagal membuat template email reset password: %w", err)
	}

	emailSender := email.NewSMTPFromEnv()
	go func() {
		if err := emailSender.SendEmail(user.Email, subject, body); err != nil {
			log.Printf("Gagal mengirim email reset password ke %s: %v", user.Email, err)
		} else {
			log.Printf("Email reset password berhasil dikirim ke %s", user.Email)
		}
	}()
	return nil
}
