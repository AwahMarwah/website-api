package auth

import (
	"net/http"
	"website-api/library/response"
	authModel "website-api/model/auth"

	"github.com/gin-gonic/gin"
)

func (c *controller) ResendVerification(ctx *gin.Context) {
	var reqBody authModel.ResendVerificationRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	statusCode, err := c.authService.ResendVerification(&reqBody)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "Verification email resent successfully", nil)
}
