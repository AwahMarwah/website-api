package content_page

import (
	"net/http"
	"website-api/library/response"
	"website-api/model/content-page"

	"github.com/gin-gonic/gin"
)

// GetBySlug
//
// @Summary		Get Content Page by Slug
// @Description Retrieve detailed information of a content page using its slug
// @Tags		2. Content Page
// @Accept		json
// @Produce		json
// @Param		slug path string true "Content Page Slug"
// @Success      200 {object} map[string]interface{} "{"data": content_page.DetailResponse, "message": ""}"
// @Failure      400 {object} map[string]interface{} "{"data": null, "message": "Bad Request"}"
// @Failure      500 {object} map[string]interface{} "{"data": null, "message": "Internal Server Error"}"
// @Router       /content-page/pages/{slug} [get]
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
