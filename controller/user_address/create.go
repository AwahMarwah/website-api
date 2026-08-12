package user_address

import (
	"net/http"
	"website-api/library/response"
	"website-api/middleware"
	userAddressModel "website-api/model/user_address"

	"github.com/gin-gonic/gin"
)

func (c *controller) Create(ctx *gin.Context) {
	var reqBody userAddressModel.CreateUserAddressRequest
	if err := ctx.ShouldBind(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	reqBody.UserID = userInfo.UserID
	if err := c.userAddressService.Create(&reqBody); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "", nil)
}
