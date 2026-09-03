package product

import (
	"net/http"
	"website-api/library/response"

	"github.com/gin-gonic/gin"
)

// @Summary Get Product Detail
// @Description Mengambil detail produk beserta variannya
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} product.ProductDetailResponse "Berhasil mengambil detail produk"
// @Router /product/{id} [get]
func (c *controller) GetProductDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	resData, err := c.productService.GetProductDetail(id)
	if err != nil {
		if err.Error() == "product not found" {
			response.Error(ctx, http.StatusNotFound, err.Error())
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "", resData)
}
