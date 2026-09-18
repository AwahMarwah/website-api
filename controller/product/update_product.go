package product

import (
	"net/http"
	"website-api/library/response"
	productModel "website-api/model/product"

	"github.com/gin-gonic/gin"
)

func (c *controller) UpdateProduct(ctx *gin.Context) {
	var reqPath productModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var req productModel.UpdateProductReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	statusCode, err := c.productService.UpdateProduct(reqPath.Id, &req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}