package order

import (
	"net/http"
	"website-api/library/response"
	"website-api/middleware"
	"website-api/model/order"

	"github.com/gin-gonic/gin"
)

// @Summary Cancel Order
// @Description Membatalkan order oleh user pemilik (hanya status PENDING)
// @Tags 4. Order
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]interface{} "Order berhasil dibatalkan"
// @Router /order/{id}/cancel [patch]
func (c *controller) CancelOrder(ctx *gin.Context) {
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

	statusCode, err := c.orderService.CancelOrder(reqPath.Id, userInfo.UserID)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}