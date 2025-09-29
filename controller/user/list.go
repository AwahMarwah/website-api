package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"website-api/library/pagination"
	"website-api/library/response"
	userModel "website-api/model/user"
)

func (c *controller) List(ctx *gin.Context) {
	var reqQuery userModel.ListUserReqQuery
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)
	resData, count, err := c.userService.List(&reqQuery)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithPage(ctx, http.StatusOK, "", resData, reqQuery.Page, reqQuery.Limit, count)
}
