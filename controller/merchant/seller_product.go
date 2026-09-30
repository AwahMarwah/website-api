package merchant

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"website-api/library/response"
	"website-api/library/pagination"
	"website-api/middleware"
	merchantModel "website-api/model/merchant"
)

// @Summary Seller List Products
// @Description Daftar produk milik seller yang sedang login
// @Tags 7. Seller Panel
// @Produce json
// @Param req query merchant.ListSellerProductReqQuery false "Query Parameters"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil daftar produk"
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /merchant/products [get]
func (c *controller) ListProducts(ctx *gin.Context) {
	var reqQuery merchantModel.ListSellerProductReqQuery
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

	resData, count, statusCode, err := c.merchantService.ListProducts(userInfo.UserID, &reqQuery)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.SuccessWithPage(ctx, statusCode, "", resData, reqQuery.Page, reqQuery.Limit, count)
}

// @Summary Seller Get Product
// @Description Detail satu produk milik seller
// @Tags 7. Seller Panel
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} merchant.SellerProductResponse "Berhasil mengambil detail produk"
// @Failure 404 {object} map[string]interface{}
// @Router /merchant/products/{id} [get]
func (c *controller) GetProduct(ctx *gin.Context) {
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

	resData, statusCode, err := c.merchantService.GetProduct(userInfo.UserID, reqPath.Id)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}

// @Summary Seller Create Product
// @Description Membuat produk baru milik seller (status draft)
// @Tags 7. Seller Panel
// @Accept json
// @Produce json
// @Param req body merchant.CreateSellerProductReq true "Body"
// @Success 201 {object} map[string]interface{} "Produk berhasil dibuat"
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /merchant/products [post]
func (c *controller) CreateProduct(ctx *gin.Context) {
	var req merchantModel.CreateSellerProductReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	productID, statusCode, err := c.merchantService.CreateProduct(userInfo.UserID, &req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", gin.H{"id": productID})
}

// @Summary Seller Update Product
// @Description Memperbarui produk milik seller
// @Tags 7. Seller Panel
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param req body merchant.UpdateSellerProductReq true "Body"
// @Success 200 {object} map[string]interface{} "Produk berhasil diperbarui"
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /merchant/products/{id} [put]
func (c *controller) UpdateProduct(ctx *gin.Context) {
	var reqPath merchantModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var req merchantModel.UpdateSellerProductReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	statusCode, err := c.merchantService.UpdateProduct(userInfo.UserID, reqPath.Id, &req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}

// @Summary Seller Delete Product
// @Description Menonaktifkan produk milik seller (soft delete)
// @Tags 7. Seller Panel
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]interface{} "Produk berhasil dinonaktifkan"
// @Failure 404 {object} map[string]interface{}
// @Router /merchant/products/{id} [delete]
func (c *controller) DeleteProduct(ctx *gin.Context) {
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

	statusCode, err := c.merchantService.DeleteProduct(userInfo.UserID, reqPath.Id)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}
