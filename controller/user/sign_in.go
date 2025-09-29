package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"website-api/common"
	"website-api/library/response"
	userModel "website-api/model/user"
)

func (c *controller) SignIn(ctx *gin.Context) {
	var reqBody userModel.SignInRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resData, statusCode, err := c.userService.SignIn(&reqBody)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, common.SuccessfullySignedIn, resData)
}
