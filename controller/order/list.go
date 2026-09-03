package order

import (
	"net/http"
	"website-api/library/pagination"
	"website-api/library/response"
	"website-api/middleware"
	"website-api/model/order"

	"github.com/gin-gonic/gin"
)

// @Summary List Orders (User)
// @Description Mengambil daftar order milik user yang login
// @Tags 4. Order
// @Produce json
// @Success 200 {object} order.OrderResponse "Berhasil mengambil daftar order"
// @Router /order [get]
func (c *controller) List(ctx *gin.Context) {
	var reqQuery order.ListOrderReqQuery
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	reqQuery.UserID = userInfo.UserID
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)

	resData, count, statusCode, err := c.orderService.List(&reqQuery)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.SuccessWithPage(ctx, statusCode, "", resData, reqQuery.Page, reqQuery.Limit, count)
}
