package auth

type (
	ForgotPasswordRequest struct {
		Email string `binding:"required,email" json:"email"`
	}

	ResetPasswordRequest struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	ResendVerificationRequest struct {
		Email string `binding:"required,email" json:"email"`
	}

	ReqHeader struct {
		Authorization string `binding:"required"`
	}

	ChangePasswordRequest struct {
		CurrentPassword string `binding:"required" json:"current_password"`
		NewPassword     string `binding:"required,min=8,max=12" json:"new_password"`
	}
)
