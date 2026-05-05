package content_page

import (
	"net/http"
	"website-api/library/pagination"
	"website-api/library/response"
	modelContentPage "website-api/model/content-page"

	"github.com/gin-gonic/gin"
)

func (c *controller) GetFaq(ctx *gin.Context) {
	var reqQuery modelContentPage.FaqListReqQuery
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)
	resData, count, err := c.contentPageService.FaqList(&reqQuery)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithPage(ctx, http.StatusOK, "", resData, reqQuery.Page, reqQuery.Limit, count)
	return
}
