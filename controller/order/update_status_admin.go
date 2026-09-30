package order

import (
	"net/http"
	"website-api/library/response"
	"website-api/middleware"
	"website-api/model/order"

	"github.com/gin-gonic/gin"
)

// @Summary Update Order Status (Admin)
// @Description Mengubah status order oleh admin dengan validasi transisi maju
// @Tags 4. Order
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param req body order.UpdateOrderStatusReq true "Body"
// @Success 200 {object} map[string]interface{} "Status order berhasil diubah"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/orders/{id}/status [patch]
func (c *controller) UpdateStatusAdmin(ctx *gin.Context) {
	var reqPath order.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var reqBody order.UpdateOrderStatusReq
	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	statusCode, err := c.orderService.UpdateStatusAdmin(reqPath.Id, userInfo.UserID, userInfo.Role, reqBody.Status, reqBody.Note)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}
