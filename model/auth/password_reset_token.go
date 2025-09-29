package auth

import "time"

type (
	PasswordResetToken struct {
		Id        string
		UserID    string
		TokenHash string
		ExpiresAt time.Time
		IsUsed    bool
		CreatedAt time.Time
	}
)
