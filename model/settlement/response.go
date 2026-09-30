package settlement

import "time"

type (
	// MerchantLedgerResponse menampilkan saldo sebagai hasil berjalan, bukan
	// hitung ulang dari nol, supaya bisa dibandingkan dengan catatan di backend.
	MerchantLedgerResponse struct {
		ID           string    `json:"id"`
		MerchantID   string    `json:"merchant_id"`
		OrderID      *string   `json:"order_id"`
		Type         string    `json:"type"`
		Amount       float64   `json:"amount"`
		BalanceAfter float64   `json:"balance_after"`
		Note         *string   `json:"note"`
		CreatedAt    time.Time `json:"created_at"`
	}

	MerchantLedgerItemResponse struct {
		ID               string  `json:"id"`
		OrderItemID      string  `json:"order_item_id"`
		Subtotal         float64 `json:"subtotal"`
		CommissionRateBP int     `json:"commission_rate_bp"`
		Commission       float64 `json:"commission"`
	}

	MerchantLedgerDetailResponse struct {
		MerchantID   string                      `json:"merchant_id"`
		Balance      float64                     `json:"balance"`
		Type         string                      `json:"type"`
		Limit        int                         `json:"limit"`
		Offset       int                         `json:"offset"`
		Items        []MerchantLedgerItemResponse `json:"items"`
	}

	PayoutResponse struct {
		ID          string     `json:"id"`
		MerchantID  string     `json:"merchant_id"`
		Amount      float64    `json:"amount"`
		Status      string     `json:"status"`
		BankAccount *string    `json:"bank_account"`
		Note        *string    `json:"note"`
		RequestedBy *string    `json:"requested_by"`
		ApprovedBy  *string    `json:"approved_by"`
		CreatedAt   time.Time  `json:"created_at"`
	}

	BalanceResponse struct {
		MerchantID      string  `json:"merchant_id"`
		Balance         float64 `json:"balance"`
		CommissionRateBP int    `json:"commission_rate_bp"`
		PendingPayout   float64 `json:"pending_payout"`
	}
)

type (
	ListLedgerReqQuery struct {
		Page   int    `form:"page"`
		Limit  int    `form:"limit"`
		Offset int    `form:"offset"`
		Type   string `form:"type"`
	}

	ListPayoutReqQuery struct {
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
		Offset     int    `form:"offset"`
		Status     string `form:"status"`
		MerchantID string `form:"merchant_id"`
	}

	CreatePayoutReq struct {
		Amount      float64 `binding:"required,gt=0" json:"amount"`
		BankAccount string  `binding:"required" json:"bank_account"`
		Note        string  `json:"note"`
	}

	ApprovePayoutReq struct {
		Note string `json:"note"`
	}

	ReqPath struct {
		Id string `uri:"id" binding:"required"`
	}
)
