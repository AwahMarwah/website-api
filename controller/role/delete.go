package role

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"website-api/common"
	"website-api/library/response"
	modelRole "website-api/model/role"
)

func (c *controller) Delete(ctx *gin.Context) {
	var reqPath modelRole.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if statusCode, err := c.roleService.Delete(&reqPath); err != nil {
		response.Error(ctx, statusCode, err.Error())
	}
	response.Success(ctx, http.StatusOK, common.SuccessfullyDeleted, nil)

}
