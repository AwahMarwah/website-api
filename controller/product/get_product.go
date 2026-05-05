package product

import (
	"net/http"
	"website-api/library/pagination"
	"website-api/library/response"
	productModel "website-api/model/product"

	"github.com/gin-gonic/gin"
)

func (c *controller) GetProduct(ctx *gin.Context) {
	var reqQuery productModel.GetListProductReqQuerry
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)
	resData, count, err := c.productService.GetProduct(&reqQuery)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithPage(ctx, http.StatusOK, "", resData, reqQuery.Page, reqQuery.Limit, count)
}
