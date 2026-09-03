package order

import (
	"net/http"
	"website-api/library/response"
	"website-api/middleware"
	"website-api/model/order"

	"github.com/gin-gonic/gin"
)

// @Summary Create Payment Link
// @Description Membuat payment link Midtrans untuk order yang masih PENDING sehingga bisa dibagikan
// @Tags 4. Order
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} order.PaymentLinkResponse "Payment link berhasil dibuat"
// @Router /order/{id}/payment-link [post]
func (c *controller) CreatePaymentLink(ctx *gin.Context) {
	var reqPath order.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	resData, statusCode, err := c.orderService.CreatePaymentLink(reqPath.Id, userInfo.UserID, userInfo.Role)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}