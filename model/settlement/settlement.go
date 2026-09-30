package settlement

import "time"

const (
	LedgerTypeCommission = "ORDER_COMMISSION"
	LedgerTypePayout     = "PAYOUT"
	LedgerTypeRefund     = "REFUND"
)

const (
	PayoutStatusPending  = "PENDING"
	PayoutStatusApproved = "APPROVED"
	PayoutStatusPaid     = "PAID"
	PayoutStatusRejected = "REJECTED"
)

type MerchantLedger struct {
	ID           string
	MerchantID   string
	OrderID      *string
	Type         string
	Amount       float64
	BalanceAfter float64
	Note         *string
	CreatedAt    time.Time
}

// MerchantLedgerItem mengunci satu item order ke satu baris komisi lewat
// order_item_id yang UNIQUE. Inilah yang membuat komisi tidak bisa tercatat
// dua kali kalau handler order dijalankan ulang.
type MerchantLedgerItem struct {
	ID                string
	LedgerID          string
	OrderItemID       string
	Subtotal          float64
	CommissionRateBP  int
	Commission        float64
	CreatedAt         time.Time
}

type Payout struct {
	ID          string
	MerchantID  string
	Amount      float64
	Status      string
	BankAccount *string
	Note        *string
	RequestedBy *string
	ApprovedBy  *string
	LedgerID    *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
