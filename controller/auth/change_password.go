package auth

import (
	"net/http"
	"website-api/library/response"
	"website-api/middleware"
	authModel "website-api/model/auth"

	"github.com/gin-gonic/gin"
)

func (c *controller) ChangePassword(ctx *gin.Context) {
	var req authModel.ChangePasswordRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	statusCode, err := c.authService.ChangePassword(userInfo.UserID, &req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "Password berhasil diubah", nil)
}
