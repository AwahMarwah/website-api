package master

import (
	"net/http"
	"website-api/library/pagination"
	"website-api/library/response"
	"website-api/model/master"

	"github.com/gin-gonic/gin"
)

// GetSubdistrict godoc
// @Summary		Get List Subdistricts
// @Description	Mengambil master list subdistricts/kelurahan berdasarkan district
// @Tags		3. Master
// @Param		req query master.GetListSubdistrictsRequest false "Query Parameters"
// @Accept		json
// @Produce		json
// @Success		200 {object} master.SwaggerSubdistrictPagination "Berhasil mengambil data subdistrict"
// @Router		/master/subdistricts [get]
func (c *controller) GetSubdistrict(ctx *gin.Context) {
	var reqQuery master.GetListSubdistrictsRequest
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)
	resData, count, err := c.masterService.GetSubdistrict(&reqQuery)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithPage(ctx, http.StatusOK, "", resData, reqQuery.Page, reqQuery.Limit, count)
}
