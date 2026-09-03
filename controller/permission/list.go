package permission

import (
	"net/http"
	"website-api/library/response"

	"github.com/gin-gonic/gin"
)

func (c *controller) List(ctx *gin.Context) {
	resData, err := c.permissionService.List()
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "", resData)
}
