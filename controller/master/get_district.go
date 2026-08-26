package master

import (
	"net/http"
	"website-api/library/pagination"
	"website-api/library/response"
	"website-api/model/master"

	"github.com/gin-gonic/gin"
)

func (c *controller) GetDistrict(ctx *gin.Context) {
	var reqQuery master.GetListDistrictsRequest
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)
	resData, count, err := c.masterService.GetDistrict(&reqQuery)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithPage(ctx, http.StatusOK, "", resData, reqQuery.Page, reqQuery.Limit, count)
}
