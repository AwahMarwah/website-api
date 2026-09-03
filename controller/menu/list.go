package menu

import (
	"net/http"
	"website-api/library/response"
	menuModel "website-api/model/menu"

	"github.com/gin-gonic/gin"
)

func (c *controller) List(ctx *gin.Context) {
	var reqQuery menuModel.ListMenuReqQuery
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resData, count, statusCode, err := c.menuService.List(reqQuery.Page, reqQuery.Limit, reqQuery.Offset)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
	_ = count
}
