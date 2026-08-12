package user_address

import (
	"net/http"
	"website-api/library/response"
	userAddressModel "website-api/model/user_address"

	"github.com/gin-gonic/gin"
)

func (c *controller) ChangePrimary(ctx *gin.Context) {
	var req userAddressModel.UpdateUserAddress
	if err := ctx.ShouldBindUri(&req.Path); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := ctx.ShouldBindJSON(&req.Body); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.userAddressService.Update(&req); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "", nil)
}
