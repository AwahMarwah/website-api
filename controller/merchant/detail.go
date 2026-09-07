package merchant

import (
	"net/http"
	"website-api/library/response"
	merchantModel "website-api/model/merchant"

	"github.com/gin-gonic/gin"
)

func (c *controller) Detail(ctx *gin.Context) {
	var reqPath merchantModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resData, statusCode, err := c.merchantService.Detail(reqPath.Id)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}