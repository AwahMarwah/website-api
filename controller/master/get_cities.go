package master

import (
	"net/http"
	"website-api/library/pagination"
	"website-api/library/response"
	"website-api/model/master"

	"github.com/gin-gonic/gin"
)

// GetCities	godoc
// @Summary		Get List Cities
// @Description	Mengambil master list cities
// @Tags		3. Master
// @Param		req query master.GetListCitiesRequest false "Query Parameters"
// @Accept		json
// @Produce		json
// @Success		200 {object} master.SwaggerCityPagination "Berhasil mengambil data city"
// @Router		/master/cities [get]
func (c *controller) GetCities(ctx *gin.Context) {
	var reqQuery master.GetListCitiesRequest
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)
	resData, count, err := c.masterService.GetCities(&reqQuery)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessWithPage(ctx, http.StatusOK, "", resData, reqQuery.Page, reqQuery.Limit, count)
}
