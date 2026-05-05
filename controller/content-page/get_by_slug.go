package content_page

import (
	"net/http"
	"website-api/library/response"
	"website-api/model/content-page"

	"github.com/gin-gonic/gin"
)

func (c *controller) GetBySlug(ctx *gin.Context) {
	var reqPath content_page.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resData, statusCode, err := c.contentPageService.Detail(&reqPath)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}
