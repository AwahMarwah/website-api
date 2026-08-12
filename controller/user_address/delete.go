package user_address

import (
	"net/http"
	"website-api/library/response"
	userAddressModel "website-api/model/user_address"

	"github.com/gin-gonic/gin"
)

func (c *controller) Delete(ctx *gin.Context) {
	var reqPath userAddressModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.userAddressService.Delete(reqPath); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "", nil)
}
