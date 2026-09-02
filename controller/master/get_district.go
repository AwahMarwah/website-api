package master

import (
	"net/http"
	"website-api/library/pagination"
	"website-api/library/response"
	"website-api/model/master"

	"github.com/gin-gonic/gin"
)

// GetDistrict	godoc
// @Summary		Get List District
// @Description	Mengambil master list district
// @Tags		3. Master
// @Param		req query master.GetListDistrictsRequest false "Query Parameters"
// @Accept		json
// @Produce		json
// @Success		200 {object} master.SwaggerDistrictPagination "Berhasil mengambil data district"
// @Router		/master/district [get]
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
