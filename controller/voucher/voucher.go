package voucher

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"website-api/database/transaction"
	"website-api/library/pagination"
	"website-api/library/response"
	"website-api/middleware"
	voucherModel "website-api/model/voucher"
	voucherRepo "website-api/repository/voucher"
	"website-api/service/voucher"
)

type controller struct {
	voucherService voucher.IService
}

func NewController(db *gorm.DB) *controller {
	return &controller{voucherService: voucher.NewService(
		voucherRepo.NewRepo(db),
		transaction.NewTransactionManager(db),
	)}
}

// @Summary Create Voucher
// @Description Membuat voucher atau kode promo
// @Tags 8. Voucher
// @Accept json
// @Produce json
// @Param req body voucher.CreateVoucherReq true "Body"
// @Success 201 {object} map[string]interface{} "Voucher berhasil dibuat"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/vouchers [post]
func (c *controller) Create(ctx *gin.Context) {
	var req voucherModel.CreateVoucherReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, statusCode, err := c.voucherService.Create(&req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", gin.H{"id": id})
}

// @Summary List Vouchers
// @Description Daftar voucher
// @Tags 8. Voucher
// @Produce json
// @Param req query voucher.ListVoucherReqQuery false "Query Parameters"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil daftar voucher"
// @Failure 500 {object} map[string]interface{}
// @Router /admin/vouchers [get]
func (c *controller) List(ctx *gin.Context) {
	var reqQuery voucherModel.ListVoucherReqQuery
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)

	resData, count, statusCode, err := c.voucherService.List(&reqQuery)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.SuccessWithPage(ctx, statusCode, "", resData, reqQuery.Page, reqQuery.Limit, count)
}

// @Summary Voucher Detail
// @Description Detail satu voucher
// @Tags 8. Voucher
// @Produce json
// @Param id path string true "Voucher ID"
// @Success 200 {object} voucher.VoucherResponse "Berhasil mengambil detail voucher"
// @Failure 404 {object} map[string]interface{}
// @Router /admin/vouchers/{id} [get]
func (c *controller) Detail(ctx *gin.Context) {
	var reqPath voucherModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resData, statusCode, err := c.voucherService.Detail(reqPath.Id)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}

// @Summary Update Voucher
// @Description Memperbarui voucher
// @Tags 8. Voucher
// @Accept json
// @Produce json
// @Param id path string true "Voucher ID"
// @Param req body voucher.UpdateVoucherReq true "Body"
// @Success 200 {object} map[string]interface{} "Voucher berhasil diperbarui"
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/vouchers/{id} [put]
func (c *controller) Update(ctx *gin.Context) {
	var reqPath voucherModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var req voucherModel.UpdateVoucherReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	statusCode, err := c.voucherService.Update(reqPath.Id, &req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}

// @Summary Delete Voucher
// @Description Menonaktifkan voucher (soft delete, histori redemption tetap tersimpan)
// @Tags 8. Voucher
// @Produce json
// @Param id path string true "Voucher ID"
// @Success 200 {object} map[string]interface{} "Voucher berhasil dinonaktifkan"
// @Failure 404 {object} map[string]interface{}
// @Router /admin/vouchers/{id} [delete]
func (c *controller) Delete(ctx *gin.Context) {
	var reqPath voucherModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	statusCode, err := c.voucherService.Delete(reqPath.Id)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}

// @Summary Validate Voucher
// @Description Mengecek voucher terhadap nilai belanja, tanpa mengubah apa pun.
// @Description Selalu mengembalikan 200 dengan flag valid; voucher tidak berlaku adalah jawaban normal.
// @Tags 8. Voucher
// @Produce json
// @Param req query voucher.ValidateVoucherReqQuery false "Query Parameters"
// @Success 200 {object} voucher.ValidateVoucherResponse "Hasil validasi voucher"
// @Router /voucher/validate [get]
func (c *controller) Validate(ctx *gin.Context) {
	var reqQuery voucherModel.ValidateVoucherReqQuery
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	resData, err := c.voucherService.Validate(reqQuery.Code, userInfo.UserID, reqQuery.Amount)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "", resData)
}
