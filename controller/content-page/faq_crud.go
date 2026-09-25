package content_page

import (
	"net/http"
	"website-api/library/pagination"
	"website-api/library/response"
	modelContentPage "website-api/model/content-page"

	"github.com/gin-gonic/gin"
)

func (c *controller) GetFaqListAdmin(ctx *gin.Context) {
	var reqQuery modelContentPage.FaqListReqQuery
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)
	resData, count, err := c.contentPageService.FaqListAdmin(&reqQuery)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithPage(ctx, http.StatusOK, "", resData, reqQuery.Page, reqQuery.Limit, count)
}

func (c *controller) CreateFaq(ctx *gin.Context) {
	var req modelContentPage.FaqCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	statusCode, err := c.contentPageService.FaqCreate(&req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "FAQ berhasil dibuat", nil)
}

func (c *controller) UpdateFaq(ctx *gin.Context) {
	var reqPath modelContentPage.FaqReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var req modelContentPage.FaqUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	statusCode, err := c.contentPageService.FaqUpdate(reqPath.Id, &req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "FAQ berhasil diperbarui", nil)
}

func (c *controller) DeleteFaq(ctx *gin.Context) {
	var reqPath modelContentPage.FaqReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	statusCode, err := c.contentPageService.FaqDelete(reqPath.Id)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "FAQ berhasil dihapus", nil)
}