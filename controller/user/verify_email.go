package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"website-api/common"
	"website-api/library/response"
	userModel "website-api/model/user"
)

func (c *controller) VerifyEmail(ctx *gin.Context) {
	var reqBody userModel.VerifyEmailRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
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
