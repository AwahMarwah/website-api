package brand

import (
	"net/http"
	"website-api/library/response"
	brandModel "website-api/model/brand"

	"github.com/gin-gonic/gin"
)

func (c *controller) Update(ctx *gin.Context) {
	var reqPath brandModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var req brandModel.UpdateBrandReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	statusCode, err := c.brandService.Update(reqPath.Id, &req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}