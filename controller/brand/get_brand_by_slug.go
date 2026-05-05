package brand

import (
	"net/http"
	"website-api/library/response"
	brandModel "website-api/model/brand"

	"github.com/gin-gonic/gin"
)

func (c *controller) GetBrandBySlug(ctx *gin.Context) {
	var reqBody brandModel.FilterBrandReq
	if err := ctx.ShouldBindUri(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
	}
	resData, err := c.brandService.GetBrandBySlug(&reqBody)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	response.Success(ctx, http.StatusOK, "", resData)
}
