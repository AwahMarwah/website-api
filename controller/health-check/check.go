package health_check

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"website-api/common"
	"website-api/library/response"
)

func (c *controller) Check(ctx *gin.Context) {
	if err := c.healthService.Check(); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, common.SuccessfullyChecked, nil)
}
