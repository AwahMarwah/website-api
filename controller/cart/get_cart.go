package cart

import (
	"net/http"
	"website-api/library/response"
	"website-api/middleware"

	"github.com/gin-gonic/gin"
)

func (c *controller) GetCart(ctx *gin.Context) {
	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	resData, statusCode, err := c.cartService.GetCartByUserID(userInfo.UserID)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
	}

	response.Success(ctx, http.StatusOK, "", resData)
}
