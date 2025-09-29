package user

type (
	ListUserReqQuery struct {
		Page   int    `form:"page"`
		Limit  int    `form:"limit"`
		Offset int    `form:"offset"`
		Search string `form:"search"`
	}
	ReqPath struct {
		Id string `uri:"id" binding:"required"`
	}
	RegisterRequest struct {
		Name        string `binding:"required,min=8" json:"name"`
		Username    string `binding:"required,min=8" json:"username"`
		PhoneNumber string `binding:"required,min=10,max=15" json:"phone_number"`
		Email       string `binding:"required,email" json:"email"`
		Password    string `binding:"omitempty,min=8,max=12" json:"password"`
		RoleName    string `json:"role_name"`
	}

	SignInRequest struct {
		Email    string `binding:"required,email" json:"email"`
		Password string `binding:"required,min=8,max=12" json:"password"`
	}

	VerifyEmailRequest struct {
		Token string `binding:"required" json:"token"`
	}
)
