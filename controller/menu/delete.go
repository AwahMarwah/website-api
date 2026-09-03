package menu

import (
	"net/http"
	"website-api/library/response"
	menuModel "website-api/model/menu"

	"github.com/gin-gonic/gin"
)

func (c *controller) Delete(ctx *gin.Context) {
	var reqPath menuModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	statusCode, err := c.menuService.Delete(reqPath.Id)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}
