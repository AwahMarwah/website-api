package role

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"website-api/library/response"
)

func (c *controller) Find(ctx *gin.Context) {
	result, err := c.roleService.Find()
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "", result)
}
