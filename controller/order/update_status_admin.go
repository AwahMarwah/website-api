package order

import (
	"net/http"
	"website-api/library/response"
	"website-api/model/order"

	"github.com/gin-gonic/gin"
)

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
	statusCode, err := c.orderService.UpdateStatusAdmin(reqPath.Id, reqBody.Status)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}