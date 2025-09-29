package role

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"website-api/common"
	"website-api/library/response"
	"website-api/middleware"
	roleModel "website-api/model/role"
)

func (c *controller) Create(ctx *gin.Context) {
	var reqBody roleModel.RoleCreateReqBody
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	fmt.Println(userInfo, "ini data usernya")
	statusCode, err := c.roleService.Create(&reqBody)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, http.StatusCreated, common.SuccessfullyCreated, nil)
}
