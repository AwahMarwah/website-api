package health_check

import (
	"net/http"
	"website-api/common"
	"website-api/library/response"

	"github.com/gin-gonic/gin"
)

// Check godoc
//
// @Summary      Health Check
// @Description  Check application health status
// @Tags         1. Health Check
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{} "{"data": null, "message": "successfully checked"}"
// @Failure      500 {object} map[string]interface{} "{"data": null, "message": "internal server error"}"
// @Router       /health [get]
func (c *controller) Check(ctx *gin.Context) {
	if err := c.healthService.Check(); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, common.SuccessfullyChecked, nil)
}
