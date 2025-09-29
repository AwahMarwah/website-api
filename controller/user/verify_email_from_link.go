package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"website-api/common"
	"website-api/library/response"
	userModel "website-api/model/user"
)

// Endpoint GET baru untuk verifikasi email dari browser
func (c *controller) VerifyEmailFromLink(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		response.Error(ctx, http.StatusBadRequest, "Token verifikasi diperlukan")
		return
	}

	reqBody := userModel.VerifyEmailRequest{
		Token: token,
	}

	err := c.userService.VerifyEmail(reqBody)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == common.TooManyVerificationAttempts {
			statusCode = http.StatusTooManyRequests
		}

		response.Error(ctx, statusCode, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, common.VerificationSuccess, nil)
}
