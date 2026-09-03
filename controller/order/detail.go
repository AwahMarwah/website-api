package order

import (
	"net/http"
	"website-api/library/response"
	"website-api/middleware"
	"website-api/model/order"

	"github.com/gin-gonic/gin"
)

// @Summary Get Order Detail
// @Description Mengambil detail order beserta item (hanya pemilik order atau admin)
// @Tags 4. Order
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} order.OrderResponse "Berhasil mengambil detail order"
// @Router /order/{id} [get]
func (c *controller) Detail(ctx *gin.Context) {
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

	resData, statusCode, err := c.orderService.Detail(reqPath.Id, userInfo.UserID, userInfo.Role)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}
