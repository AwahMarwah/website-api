package user

import (
	"net/http"
	"website-api/common"
	"website-api/library/response"

	"github.com/gin-gonic/gin"
)

func (c *controller) SignOut(ctx *gin.Context) {
	if err := c.userService.SignOut(ctx.MustGet("user_id").(string)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, common.SuccessfullySignedOut, nil)
}
