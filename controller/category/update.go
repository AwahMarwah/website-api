package category

import (
	"net/http"
	"website-api/library/response"
	categoryModel "website-api/model/category"

	"github.com/gin-gonic/gin"
)

func (c *controller) Update(ctx *gin.Context) {
	var reqPath categoryModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var req categoryModel.UpdateCategoryReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	statusCode, err := c.categoryService.Update(reqPath.Id, &req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}