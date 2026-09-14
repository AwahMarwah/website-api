package product

import (
	"website-api/library/response"
	"github.com/gin-gonic/gin"
)

func (c *controller) DeleteImage(ctx *gin.Context) {
	productID := ctx.Param("id")
	imageID := ctx.Param("imageId")
	statusCode, err := c.productService.DeleteImage(productID, imageID)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}
