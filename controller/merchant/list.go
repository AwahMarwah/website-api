package merchant

import (
	"website-api/library/response"

	"github.com/gin-gonic/gin"
)

// @Summary Get Merchants
// @Description Mengambil daftar merchant (toko) yang aktif
// @Tags Merchant
// @Accept json
// @Produce json
// @Success 200 {object} merchant.MerchantResponse "Berhasil mengambil daftar merchant"
// @Router /merchant [get]
func (c *controller) List(ctx *gin.Context) {
	resData, statusCode, err := c.merchantService.List()
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}