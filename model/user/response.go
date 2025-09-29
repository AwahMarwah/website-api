package user

import (
	"time"
	"website-api/model/role"
)

type (
	DetailResponse struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		UserName    string    `json:"user_name"`
		Email       string    `json:"email"`
		PhoneNumber string    `json:"phone_number"`
		IsVerified  bool      `json:"is_verified"`
		CreatedAt   string    `json:"created_at"`
		UpdatedAt   string    `json:"updated_at"`
		Role        role.Role `json:"role"`
	}

	ListUserResponse struct {
		Id                 string     `json:"id"`
		Name               string     `json:"name"`
		UserName           string     `json:"user_name"`
		Picture            string     `json:"picture"`
		Email              string     `json:"email"`
		PhoneNumber        string     `json:"phone_number"`
		IsVerified         bool       `json:"is_verified"`
		CreatedAt          time.Time  `json:"created_at"`
		CreatedAtFormatted string     `json:"created_at_formatted"`
		CreatedBy          string     `json:"created_by"`
		UpdatedAt          *time.Time `json:"updated_at"`
		UpdatedByFormatted string     `json:"updated_at_formatted"`
		UpdatedBy          string     `json:"updated_by"`
		DeletedAt          *time.Time `json:"deleted_at"`
		DeletedByFormatted string     `json:"deleted_by_formatted"`
		DeletedBy          string     `json:"deleted_by"`
	}
	SignInResponse struct {
		Token     string      `json:"token"`
		User      UserProfile `json:"user"`
		ExpiresAt int64       `json:"expires_at"`
	}

	UserProfile struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Username    string    `json:"username"`
		Email       string    `json:"email"`
		PhoneNumber string    `json:"phone_number"`
		IsVerified  bool      `json:"is_verified"`
		Role        role.Role `json:"role"`
	}
)
