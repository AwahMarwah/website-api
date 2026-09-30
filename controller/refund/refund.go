package refund

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"website-api/database/transaction"
	"website-api/library/pagination"
	"website-api/library/response"
	"website-api/middleware"
	refundModel "website-api/model/refund"
	orderRepo "website-api/repository/order"
	refundRepo "website-api/repository/refund"
	settlementRepo "website-api/repository/settlement"
	settlementService "website-api/service/settlement"
	refundService "website-api/service/refund"
	midtransProvider "website-api/third-party/provider/midtrans"
)

type controller struct {
	refundService refundService.IService
}

func NewController(db *gorm.DB) *controller {
	settlement := settlementService.NewService(settlementRepo.NewRepo(db))
	return &controller{refundService: refundService.NewService(
		refundRepo.NewRepo(db),
		orderRepo.NewRepo(db),
		settlementRepo.NewRepo(db),
		settlement,
		transaction.NewTransactionManager(db),
		midtransProvider.NewClient(),
	)}
}

// @Summary Request Refund
// @Description MengAJUKAN refund untuk order milik sendiri. Dana belum berpindah sampai disetujui finance.
// @Tags 10. Refund
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param req body refund.CreateRefundReq true "Body"
// @Success 201 {object} map[string]interface{} "Permintaan refund berhasil dibuat"
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /order/{id}/refund [post]
func (c *controller) Request(ctx *gin.Context) {
	var reqPath struct {
		Id string `uri:"id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var req refundModel.CreateRefundReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	refundID, statusCode, err := c.refundService.Request(reqPath.Id, userInfo.UserID, &req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", gin.H{"id": refundID})
}

// @Summary List Refunds
// @Description Daftar permintaan refund
// @Tags 10. Refund
// @Produce json
// @Param req query refund.ListRefundReqQuery false "Query Parameters"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil daftar refund"
// @Failure 500 {object} map[string]interface{}
// @Router /admin/refunds [get]
func (c *controller) List(ctx *gin.Context) {
	var reqQuery refundModel.ListRefundReqQuery
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)

	resData, count, statusCode, err := c.refundService.List(&reqQuery)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.SuccessWithPage(ctx, statusCode, "", resData, reqQuery.Page, reqQuery.Limit, count)
}

// @Summary Refund Detail
// @Description Detail satu permintaan refund
// @Tags 10. Refund
// @Produce json
// @Param id path string true "Refund ID"
// @Success 200 {object} refund.RefundResponse "Berhasil mengambil detail refund"
// @Failure 404 {object} map[string]interface{}
// @Router /admin/refunds/{id} [get]
func (c *controller) Detail(ctx *gin.Context) {
	var reqPath refundModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resData, statusCode, err := c.refundService.Detail(reqPath.Id)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}

// @Summary Approve Refund
// @Description Menyetujui permintaan refund. Dana dikembalikan saat endpoint process dipanggil.
// @Tags 10. Refund
// @Accept json
// @Produce json
// @Param id path string true "Refund ID"
// @Success 200 {object} map[string]interface{} "Refund berhasil disetujui"
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /admin/refunds/{id}/approve [patch]
func (c *controller) Approve(ctx *gin.Context) {
	var reqPath refundModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	statusCode, err := c.refundService.Approve(reqPath.Id, userInfo.UserID)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}

// @Summary Reject Refund
// @Description Menolak permintaan refund
// @Tags 10. Refund
// @Accept json
// @Produce json
// @Param id path string true "Refund ID"
// @Param req body refund.ApproveRefundReq true "Body"
// @Success 200 {object} map[string]interface{} "Refund berhasil ditolak"
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /admin/refunds/{id}/reject [patch]
func (c *controller) Reject(ctx *gin.Context) {
	var reqPath refundModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var req refundModel.ApproveRefundReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	statusCode, err := c.refundService.Reject(reqPath.Id, userInfo.UserID, req.Note)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}

// @Summary Process Refund
// @Description Mengirim dana kembali ke pembeli lewat payment provider, mengembalikan stok,
// @Description dan mengoreksi saldo seller.
// @Tags 10. Refund
// @Produce json
// @Param id path string true "Refund ID"
// @Success 200 {object} map[string]interface{} "Refund berhasil diproses"
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 502 {object} map[string]interface{}
// @Router /admin/refunds/{id}/process [post]
func (c *controller) Process(ctx *gin.Context) {
	var reqPath refundModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	statusCode, err := c.refundService.Process(reqPath.Id)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}
