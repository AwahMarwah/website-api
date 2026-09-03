package menu

import (
	"net/http"
	"website-api/library/response"
	menuModel "website-api/model/menu"

	"github.com/gin-gonic/gin"
)

func (c *controller) Update(ctx *gin.Context) {
	var reqPath menuModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var reqBody menuModel.MenuUpdateReq
	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resData, statusCode, err := c.menuService.Update(reqPath.Id, &reqBody)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}
