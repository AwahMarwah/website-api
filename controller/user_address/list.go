package user_address

import (
	"net/http"
	"website-api/library/response"
	"website-api/middleware"

	"github.com/gin-gonic/gin"
)

func (c *controller) List(ctx *gin.Context) {
	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	resData, err := c.userAddressService.Find(userInfo.UserID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "", resData)
}
