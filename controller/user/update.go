package user

import (
	"net/http"
	"website-api/library/response"
	modelUser "website-api/model/user"

	"github.com/gin-gonic/gin"
)

func (c *controller) Update(ctx *gin.Context) {
	var req modelUser.UserUpdateRequest
	if err := ctx.ShouldBindUri(&req.Path); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := ctx.ShouldBindJSON(&req.Body); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	statusCode, err := c.userService.Update(&req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
	}
	response.Success(ctx, statusCode, "", nil)
}
