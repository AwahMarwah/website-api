package order

import (
	"net/http"
	"website-api/library/pagination"
	"website-api/library/response"
	"website-api/model/order"

	"github.com/gin-gonic/gin"
)

// @Summary List Orders (Admin)
// @Description Mengambil semua order (khusus admin)
// @Tags 4. Order
// @Produce json
// @Param status query string false "Filter status order"
// @Success 200 {object} order.OrderResponse "Berhasil mengambil daftar order"
// @Router /admin/orders [get]
func (c *controller) ListAdmin(ctx *gin.Context) {
	var reqQuery order.ListOrderReqQuery
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)

	resData, count, statusCode, err := c.orderService.ListAdmin(&reqQuery)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.SuccessWithPage(ctx, statusCode, "", resData, reqQuery.Page, reqQuery.Limit, count)
}
