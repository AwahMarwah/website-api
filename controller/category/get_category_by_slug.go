package category

import (
	"net/http"
	"website-api/library/response"
	categoryModel "website-api/model/category"

	"github.com/gin-gonic/gin"
)

func (c *controller) GetCategoryBySlug(ctx *gin.Context) {
	var reqPath categoryModel.CategorySlugPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resData, statusCode, err := c.categoryService.GetCategoryBySlug(reqPath.Slug)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}
