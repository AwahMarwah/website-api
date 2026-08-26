package master

import (
	"net/http"
	"website-api/library/pagination"
	"website-api/library/response"
	masterModel "website-api/model/master"

	"github.com/gin-gonic/gin"
)

// GetProvince godoc
// @Summary		Get List Province
// @Description	Mengambil master list province
// @Tags		province
// @Accept		json
// @Produce		json
// @Success		200 {object}  map[string]interface{}
// @Router		/master/province [get]
func (c *controller) GetProvince(ctx *gin.Context) {
	var reqQuery masterModel.GetListProvinceRequest
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)
	resData, count, err := c.masterService.GetProvince(&reqQuery)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithPage(ctx, http.StatusOK, "", resData, reqQuery.Page, reqQuery.Limit, count)
}
