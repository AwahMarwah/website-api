package brand

import (
	"net/http"
	"website-api/library/response"
	brandModel "website-api/model/brand"

	"github.com/gin-gonic/gin"
)

func (c *controller) Create(ctx *gin.Context) {
	var req brandModel.CreateBrandReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	statusCode, err := c.brandService.Create(&req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}