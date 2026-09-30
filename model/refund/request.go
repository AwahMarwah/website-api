package refund

import "time"

type (
	CreateRefundReq struct {
		Amount float64 `binding:"required,gt=0" json:"amount"`
		Reason string  `binding:"required" json:"reason"`
	}

	ApproveRefundReq struct {
		Note string `json:"note"`
	}

	ReqPath struct {
		Id string `uri:"id" binding:"required"`
	}

	ListRefundReqQuery struct {
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
		Offset     int    `form:"offset"`
		Status     string `form:"status"`
		OrderID    string `form:"order_id"`
		MerchantID string `form:"merchant_id"`
	}
)

type (
	RefundResponse struct {
		ID               string     `json:"id"`
		OrderID          string     `json:"order_id"`
		MerchantID       *string    `json:"merchant_id"`
		RequestedBy      string     `json:"requested_by"`
		Amount           float64    `json:"amount"`
		Reason           *string    `json:"reason"`
		Status           string     `json:"status"`
		Provider         string     `json:"provider"`
		ProviderRefundID *string    `json:"provider_refund_id"`
		ApprovedBy       *string    `json:"approved_by"`
		Restocked        bool       `json:"restocked"`
		FailureReason    *string    `json:"failure_reason"`
		CreatedAt        time.Time  `json:"created_at"`
		UpdatedAt        *time.Time `json:"updated_at"`
	}
)
