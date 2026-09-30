package settlement

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"website-api/library/pagination"
	"website-api/library/response"
	"website-api/middleware"
	settlementModel "website-api/model/settlement"
)

// @Summary List Payouts (Finance)
// @Description Daftar pengajuan pencairan dana
// @Tags 9. Settlement
// @Produce json
// @Param req query settlement.ListPayoutReqQuery false "Query Parameters"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil daftar payout"
// @Failure 500 {object} map[string]interface{}
// @Router /admin/payouts [get]
func (c *controller) ListPayouts(ctx *gin.Context) {
	var reqQuery settlementModel.ListPayoutReqQuery
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)

	resData, count, statusCode, err := c.payoutService.ListPayouts(&reqQuery)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.SuccessWithPage(ctx, statusCode, "", resData, reqQuery.Page, reqQuery.Limit, count)
}

// @Summary Approve Payout
// @Description Menyetujui payout dan mencatat pencairan ke ledger seller
// @Tags 9. Settlement
// @Accept json
// @Produce json
// @Param id path string true "Payout ID"
// @Param req body settlement.ApprovePayoutReq true "Body"
// @Success 200 {object} map[string]interface{} "Payout berhasil disetujui"
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/payouts/{id}/approve [patch]
func (c *controller) ApprovePayout(ctx *gin.Context) {
	var reqPath settlementModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var req settlementModel.ApprovePayoutReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	statusCode, err := c.payoutService.ApprovePayout(reqPath.Id, userInfo.UserID, &req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}

// @Summary Mark Payout Paid
// @Description Menandai payout yang sudah disetujui sebagai lunas ditransfer
// @Tags 9. Settlement
// @Produce json
// @Param id path string true "Payout ID"
// @Success 200 {object} map[string]interface{} "Payout ditandai lunas"
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /admin/payouts/{id}/paid [patch]
func (c *controller) MarkPaid(ctx *gin.Context) {
	var reqPath settlementModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	statusCode, err := c.payoutService.MarkPaid(reqPath.Id, userInfo.UserID)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}

// @Summary Reject Payout
// @Description Menolak pengajuan payout. Saldo seller tidak terpengaruh.
// @Tags 9. Settlement
// @Accept json
// @Produce json
// @Param id path string true "Payout ID"
// @Param req body settlement.ApprovePayoutReq true "Body"
// @Success 200 {object} map[string]interface{} "Payout berhasil ditolak"
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /admin/payouts/{id}/reject [patch]
func (c *controller) RejectPayout(ctx *gin.Context) {
	var reqPath settlementModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var req settlementModel.ApprovePayoutReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	statusCode, err := c.payoutService.RejectPayout(reqPath.Id, userInfo.UserID, req.Note)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}
