package role

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"website-api/library/response"
	"website-api/model/role"
)

func (c *controller) Detail(ctx *gin.Context) {
	var reqPath role.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resData, statusCode, err := c.roleService.Detail(&reqPath)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}
