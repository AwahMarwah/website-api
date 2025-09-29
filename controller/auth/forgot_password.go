package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"website-api/library/response"
	authModel "website-api/model/auth"
)

func (c *controller) ForgotPassword(ctx *gin.Context) {
	var reqBody authModel.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	statusCode, message, err := c.authService.ForgotPassword(&reqBody)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}

	response.Success(ctx, statusCode, message, nil)
}
