package brand

import (
	"net/http"
	"website-api/library/response"
	brandModel "website-api/model/brand"

	"github.com/gin-gonic/gin"
)

func (c *controller) Delete(ctx *gin.Context) {
	var reqPath brandModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	statusCode, err := c.brandService.Delete(reqPath.Id)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}