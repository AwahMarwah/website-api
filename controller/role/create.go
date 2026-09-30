package role

import (
	"net/http"
	"website-api/common"
	"website-api/library/response"
	"website-api/middleware"
	roleModel "website-api/model/role"

	"github.com/gin-gonic/gin"
)

func (c *controller) Create(ctx *gin.Context) {
	var reqBody roleModel.RoleCreateReqBody
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	// Context user dipakai hanya untuk otorisasi; middleware sudah menolak non-admin.
	if _, err := middleware.GetUserFromContext(ctx); err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	statusCode, err := c.roleService.Create(&reqBody)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, http.StatusCreated, common.SuccessfullyCreated, nil)
}
