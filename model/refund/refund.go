package refund

import "time"

const (
	StatusPending    = "PENDING"
	StatusApproved   = "APPROVED"
	StatusProcessing = "PROCESSING"
	StatusCompleted  = "COMPLETED"
	StatusRejected   = "REJECTED"
	StatusFailed     = "FAILED"
)

type Refund struct {
	ID                string
	OrderID           string
	MerchantID        *string
	RequestedBy       string
	Amount            float64
	Reason            *string
	Status            string
	Provider          string
	ProviderRefundID  *string
	ApprovedBy        *string
	Restocked         bool
	FailureReason     *string
	CreatedAt         time.Time
	UpdatedAt         *time.Time
}
