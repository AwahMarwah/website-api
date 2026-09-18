package product

import (
	"net/http"
	"website-api/library/response"
	productModel "website-api/model/product"

	"github.com/gin-gonic/gin"
)

func (c *controller) CreateProduct(ctx *gin.Context) {
	var req productModel.CreateProductReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	statusCode, err := c.productService.CreateProduct(&req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}