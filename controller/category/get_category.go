package category

import (
	"net/http"
	"website-api/library/pagination"
	"website-api/library/response"
	"website-api/model/category"

	"github.com/gin-gonic/gin"
)

func (c *controller) GetCategory(ctx *gin.Context) {
	var reqQuery category.FilterCategory
	if err := ctx.ShouldBind(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
	}
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)
	resData, count, err := c.categoryService.GetCategory(&reqQuery)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	response.SuccessWithPage(ctx, http.StatusOK, "", resData, reqQuery.Limit, reqQuery.Offset, count)
}
