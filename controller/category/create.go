package category

import (
	"net/http"
	"website-api/library/response"
	categoryModel "website-api/model/category"

	"github.com/gin-gonic/gin"
)

func (c *controller) Create(ctx *gin.Context) {
	var req categoryModel.CreateCategoryReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	statusCode, err := c.categoryService.Create(&req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}