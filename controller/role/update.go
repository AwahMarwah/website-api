package role

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"website-api/common"
	"website-api/library/response"
	"website-api/model/role"
)

func (c *controller) Update(ctx *gin.Context) {
	var req role.UpdateReq
	if err := ctx.ShouldBindUri(&req.Path); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := ctx.ShouldBindJSON(&req.Body); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if statusCode, err := c.roleService.Update(&req); err != nil {
		response.Error(ctx, statusCode, err.Error())
	}
	response.Success(ctx, http.StatusOK, common.SuccessfullyUpdated, nil)
}
