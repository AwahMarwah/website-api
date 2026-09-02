package order

import (
	"net/http"
	"website-api/library/response"
	"website-api/middleware"
	"website-api/model/order"

	"github.com/gin-gonic/gin"
)

// @Summary Checkout
// @Description Membuat order baru dan menginisiasi pembayaran via Midtrans Snap
// @Tags 4. Order
// @Accept json
// @Produce json
// @Param req body order.CheckoutRequest true "Checkout Request"
// @Success 200 {object} order.CheckoutResponse "Order berhasil dibuat"
// @Router /order [post]
func (c *controller) Checkout(ctx *gin.Context) {
	var reqBody order.CheckoutRequest
	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	reqBody.UserID = userInfo.UserID
	resData, statusCode, err := c.orderService.Checkout(&reqBody)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}
