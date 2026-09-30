package settlement

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"website-api/common"
	"website-api/database/transaction"
	"website-api/library/pagination"
	"website-api/library/response"
	"website-api/middleware"
	settlementModel "website-api/model/settlement"
	merchantRepo "website-api/repository/merchant"
	settlementRepo "website-api/repository/settlement"
	settlementService "website-api/service/settlement"
)

type controller struct {
	settlementService settlementService.IService
	payoutService     settlementService.IPayoutService
	merchantRepo      merchantRepo.IRepo
}

func NewController(db *gorm.DB) *controller {
	repo := settlementRepo.NewRepo(db)
	return &controller{
		settlementService: settlementService.NewService(repo),
		payoutService:     settlementService.NewPayoutService(repo, transaction.NewTransactionManager(db)),
		merchantRepo:      merchantRepo.NewRepo(db),
	}
}

// merchantIDFor menolerolkan userID menjadi merchant aktif. Semua endpoint seller
// memakai ini supaya merchant_id tidak pernah diambil dari request.
func (c *controller) merchantIDFor(userID string) (string, int, error) {
	if userID == "" {
		return "", http.StatusUnauthorized, fmt.Errorf("user ID tidak ditemukan dalam context")
	}

	m, err := c.merchantRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || m.ID == "" {
			return "", http.StatusForbidden, fmt.Errorf("%s", common.MerchantNotRegistered)
		}
		return "", http.StatusInternalServerError, fmt.Errorf("gagal mengambil data merchant: %w", err)
	}
	if !m.IsActive {
		return "", http.StatusForbidden, fmt.Errorf("%s", common.MerchantNotApproved)
	}
	return m.ID, http.StatusOK, nil
}

// @Summary Merchant Ledger
// @Description Daftar catatan keuangan milik seller (komisi, payout, refund)
// @Tags 9. Settlement
// @Produce json
// @Param req query settlement.ListLedgerReqQuery false "Query Parameters"
// @Success 200 {object} map[string]interface{} "Berhasil mengambil ledger"
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /merchant/ledger [get]
func (c *controller) ListLedgers(ctx *gin.Context) {
	var reqQuery settlementModel.ListLedgerReqQuery
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

	merchantID, statusCode, err := c.merchantIDFor(userInfo.UserID)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}

	resData, count, statusCode, err := c.settlementService.ListLedgers(merchantID, &reqQuery)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.SuccessWithPage(ctx, statusCode, "", resData, reqQuery.Page, reqQuery.Limit, count)
}

// @Summary Merchant Balance
// @Description Saldo seller saat ini setelah dikurangi komisi dan payout
// @Tags 9. Settlement
// @Produce json
// @Success 200 {object} settlement.BalanceResponse "Berhasil mengambil saldo"
// @Failure 403 {object} map[string]interface{}
// @Router /merchant/balance [get]
func (c *controller) Balance(ctx *gin.Context) {
	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	merchantID, statusCode, err := c.merchantIDFor(userInfo.UserID)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}

	resData, statusCode, err := c.settlementService.Balance(merchantID)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}

// @Summary Request Payout
// @Description Seller mengajukan pencairan dana. Saldo baru berkurang setelah disetujui finance.
// @Tags 9. Settlement
// @Accept json
// @Produce json
// @Param req body settlement.CreatePayoutReq true "Body"
// @Success 201 {object} map[string]interface{} "Pengajuan payout berhasil dibuat"
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /merchant/payouts [post]
func (c *controller) RequestPayout(ctx *gin.Context) {
	var req settlementModel.CreatePayoutReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	merchantID, statusCode, err := c.merchantIDFor(userInfo.UserID)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}

	payoutID, statusCode, err := c.payoutService.RequestPayout(merchantID, userInfo.UserID, &req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", gin.H{"id": payoutID})
}
