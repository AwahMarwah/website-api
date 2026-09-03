package menu

import (
	"net/http"
	"website-api/library/response"
	"website-api/middleware"

	"github.com/gin-gonic/gin"
)

func (c *controller) GetMyMenus(ctx *gin.Context) {
	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	resData, statusCode, err := c.menuService.GetMyMenus(userInfo.Role)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}
