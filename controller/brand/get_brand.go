package brand

import (
	"net/http"
	"website-api/library/pagination"
	"website-api/library/response"
	brandModel "website-api/model/brand"

	"github.com/gin-gonic/gin"
)

func (c *controller) GetBrand(ctx *gin.Context) {
	var reqQuery brandModel.BrandReqQuery
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
	}
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)
	resData, count, err := c.brandService.GetBrand(&reqQuery)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessWithPage(ctx, http.StatusOK, "", resData, reqQuery.Page, reqQuery.Limit, count)
}
