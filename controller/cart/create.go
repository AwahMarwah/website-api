package cart

import (
	"net/http"
	"website-api/common"
	"website-api/library/response"
	"website-api/middleware"
	cartModel "website-api/model/cart"

	"github.com/gin-gonic/gin"
)

func (c *controller) Create(ctx *gin.Context) {
	var reqBody cartModel.CartRequestBody
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	reqBody.UserID = userInfo.UserID

	statusCode, err := c.cartService.Create(&reqBody)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, common.SuccessfullyCreated, nil)

}
