package order

import (
	"net/http"
	"website-api/library/response"
	"website-api/middleware"
	"website-api/model/order"

	"github.com/gin-gonic/gin"
)

func (c *controller) Checkout(ctx *gin.Context) {
	var reqBody order.CheckoutRequest
	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
	}

	reqBody.UserID = userInfo.UserID
	if err := c.orderService.Checkout(&reqBody); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "", nil)
}
