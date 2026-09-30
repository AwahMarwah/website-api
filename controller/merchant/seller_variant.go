package merchant

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"website-api/library/response"
	"website-api/middleware"
	merchantModel "website-api/model/merchant"
)

// @Summary Seller List Variants
// @Description Daftar variant milik satu produk milik seller
// @Tags 7. Seller Panel
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil daftar variant"
// @Failure 404 {object} map[string]interface{}
// @Router /merchant/products/{id}/variants [get]
func (c *controller) ListVariants(ctx *gin.Context) {
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

	resData, statusCode, err := c.merchantService.ListVariants(userInfo.UserID, reqPath.Id)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}

// @Summary Seller Create Variant
// @Description Menambah variant baru ke produk milik seller
// @Tags 7. Seller Panel
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param req body merchant.CreateSellerVariantReq true "Body"
// @Success 201 {object} map[string]interface{} "Variant berhasil dibuat"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /merchant/products/{id}/variants [post]
func (c *controller) CreateVariant(ctx *gin.Context) {
	var reqPath merchantModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var req merchantModel.CreateSellerVariantReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	variantID, statusCode, err := c.merchantService.CreateVariant(userInfo.UserID, reqPath.Id, &req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", gin.H{"id": variantID})
}

// @Summary Seller Update Variant
// @Description Memperbarui variant milik seller. Stok hanya bisa ditambah, tidak dikurangi manual.
// @Tags 7. Seller Panel
// @Accept json
// @Produce json
// @Param variantId path string true "Variant ID"
// @Param req body merchant.UpdateSellerVariantReq true "Body"
// @Success 200 {object} map[string]interface{} "Variant berhasil diperbarui"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /merchant/variants/{variantId} [put]
func (c *controller) UpdateVariant(ctx *gin.Context) {
	var reqPath struct {
		VariantId string `uri:"variantId" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var req merchantModel.UpdateSellerVariantReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	statusCode, err := c.merchantService.UpdateVariant(userInfo.UserID, reqPath.VariantId, &req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}

// @Summary Seller Delete Variant
// @Description Menonaktifkan variant milik seller (soft delete)
// @Tags 7. Seller Panel
// @Produce json
// @Param variantId path string true "Variant ID"
// @Success 200 {object} map[string]interface{} "Variant berhasil dinonaktifkan"
// @Failure 404 {object} map[string]interface{}
// @Router /merchant/variants/{variantId} [delete]
func (c *controller) DeleteVariant(ctx *gin.Context) {
	var reqPath struct {
		VariantId string `uri:"variantId" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	statusCode, err := c.merchantService.DeleteVariant(userInfo.UserID, reqPath.VariantId)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}
