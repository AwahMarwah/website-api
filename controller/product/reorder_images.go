package product

import (
	"net/http"
	"website-api/library/response"
	productModel "website-api/model/product"

	"github.com/gin-gonic/gin"
)

type reorderRequest struct {
	Images []productModel.ImageSortOrder `json:"images" binding:"required,dive"`
}

func (c *controller) ReorderImages(ctx *gin.Context) {
	productID := ctx.Param("id")

	var reqBody reorderRequest
	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	statusCode, err := c.productService.ReorderImages(productID, reqBody.Images)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}
