package auth

type (
	ForgotPasswordRequest struct {
		Email string `binding:"required,email" json:"email"`
	}

	ResetPasswordRequest struct {
		Token       string `json:""token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	ResendVerificationRequest struct {
		Email string `binding:"required,email" json:"email"`
	}
)
