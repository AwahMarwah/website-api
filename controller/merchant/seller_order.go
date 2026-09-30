package merchant

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"website-api/library/pagination"
	"website-api/library/response"
	"website-api/middleware"
	merchantModel "website-api/model/merchant"
)

// @Summary Seller List Orders
// @Description Daftar order yang memuat barang milik seller yang sedang login
// @Tags 7. Seller Panel
// @Produce json
// @Param req query merchant.ListSellerOrderReqQuery false "Query Parameters"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil daftar order"
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /merchant/orders [get]
func (c *controller) ListOrders(ctx *gin.Context) {
	var reqQuery merchantModel.ListSellerOrderReqQuery
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	resData, count, statusCode, err := c.merchantService.ListOrders(userInfo.UserID, &reqQuery)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.SuccessWithPage(ctx, statusCode, "", resData, reqQuery.Page, reqQuery.Limit, count)
}

// @Summary Seller Get Order
// @Description Detail order, dipangkas hanya ke item milik seller yang requesting
// @Tags 7. Seller Panel
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} merchant.SellerOrderResponse "Berhasil mengambil detail order"
// @Failure 404 {object} map[string]interface{}
// @Router /merchant/orders/{id} [get]
func (c *controller) GetOrder(ctx *gin.Context) {
	var reqPath merchantModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	resData, statusCode, err := c.merchantService.GetOrder(userInfo.UserID, reqPath.Id)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}

// @Summary Seller Dashboard Stats
// @Description Ringkasan dashboard seller: produk, order, pendapatan, dan saldo
// @Tags 7. Seller Panel
// @Produce json
// @Success 200 {object} merchant.SellerStatsResponse "Berhasil mengambil ringkasan"
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /merchant/stats [get]
func (c *controller) Stats(ctx *gin.Context) {
	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	resData, statusCode, err := c.merchantService.Stats(userInfo.UserID)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}
