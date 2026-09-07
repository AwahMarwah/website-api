package shipping

import (
	"net/http"
	"website-api/library/response"
	"website-api/middleware"
	shippingModel "website-api/model/shipping"

	"github.com/gin-gonic/gin"
)

// @Summary Quote Shipping Cost
// @Description Menghitung ongkir real per merchant untuk item checkout
// @Tags Shipping
// @Accept json
// @Produce json
// @Param req body shipping.ShippingCostRequest true "Quote Request"
// @Success 200 {object} shipping.MerchantShipping "Berhasil menghitung ongkir"
// @Router /shipping/cost [post]
func (c *controller) Cost(ctx *gin.Context) {
	var reqBody shippingModel.ShippingCostRequest
	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	resData, statusCode, err := c.shippingService.Quote(&reqBody, userInfo.UserID)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}