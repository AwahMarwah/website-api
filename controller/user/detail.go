package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"website-api/library/response"
	userModel "website-api/model/user"
)

func (c *controller) Detail(ctx *gin.Context) {
	var reqPath userModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resData, statusCode, err := c.userService.Detail(&reqPath)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}
