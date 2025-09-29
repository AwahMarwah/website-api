package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"website-api/common"
	"website-api/library/response"
	userModel "website-api/model/user"
)

func (c *controller) SignUp(ctx *gin.Context) {
	var reqBody userModel.RegisterRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	// check reqBody password
	if reqBody.Password == "" {
		reqBody.Password = common.DefaultPassword
	}

	statusCode, err := c.userService.SignUp(&reqBody)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}

	response.Success(ctx, http.StatusCreated, "User successfully registered. Please check your email for verification.", nil)
}
