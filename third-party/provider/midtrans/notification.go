package midtrans

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"website-api/common"
)

// NotificationPayload adalah struktur payload yang dikirim Midtrans ke webhook
type NotificationPayload struct {
	TransactionTime     string           `json:"transaction_time"`
	TransactionStatus   string           `json:"transaction_status"`
	TransactionID       string           `json:"transaction_id"`
	StatusCode          string           `json:"status_code"`
	SignatureKey        string           `json:"signature_key"`
	OrderID             string           `json:"order_id"`
	GrossAmount         string           `json:"gross_amount"`
	PaymentType         string           `json:"payment_type"`
	FraudStatus         string           `json:"fraud_status"`
	Currency            string           `json:"currency"`
	SettlementTime      string           `json:"settlement_time"`
	MerchantID          string           `json:"merchant_id"`
	VaNumbers           []VaNumber       `json:"va_numbers"`
	Store               string           `json:"store"`
	PaymentCode         string           `json:"payment_code"`
	ExpiryTime          string           `json:"expiry_time"`
	StatusMessage       string           `json:"status_message"`
}

type VaNumber struct {
	Bank     string `json:"bank"`
	VaNumber string `json:"va_number"`
}

// VerifySignature memverifikasi signature key dari notifikasi Midtrans.
// Signature dihitung dari SHA512(order_id + status_code + gross_amount + server_key).
func (c *Client) VerifySignature(payload NotificationPayload) bool {
	plain := fmt.Sprintf(
		"%s%s%s%s",
		payload.OrderID,
		payload.StatusCode,
		payload.GrossAmount,
		c.serverKey,
	)
	hashed := sha512.Sum512([]byte(plain))
	expected := hex.EncodeToString(hashed[:])
	return expected == payload.SignatureKey
}

// ResolveStatus memetakan status transaksi Midtrans ke status order internal.
//
// Status refund TIDAK lagi dipetakan ke CANCELLED: uang yang sudah keluar untuk order
// yang sama sekali berbeda dari order yang dibatalkan sebelum dibayar. Per refunded
// ditangani sebagai status tersendiri supaya tidak hilang jejaknya.
func (p NotificationPayload) ResolveStatus() string {
	switch p.TransactionStatus {
	case "capture":
		if p.FraudStatus == "accept" {
			return common.OrderStatusPaid
		}
		// fraud status challenge/deny dianggap pending untuk review
		return common.OrderStatusPending
	case "settlement":
		return common.OrderStatusPaid
	case "pending":
		return common.OrderStatusPending
	case "deny", "cancel", "failure":
		return common.OrderStatusCancelled
	case "expire":
		return common.OrderStatusExpired
	case "refund", "partial_refund":
		return common.OrderStatusRefunded
	default:
		return common.OrderStatusPending
	}
}

// IsRefund menandai notifikasi yang resulted dalam pengembalian dana.
func (p NotificationPayload) IsRefund() bool {
	return p.TransactionStatus == "refund" || p.TransactionStatus == "partial_refund"
}

// IsFinal menandai notifikasi yang tidak akan memicu perubahan status lagi.
// Midtrans mengirim beberapa notifikasi beruntun untuk satu transaksi; tanpa ini
// email konfirmasi bisa terkirim berkali-kali.
func (p NotificationPayload) IsFinal() bool {
	switch p.TransactionStatus {
	case "settlement", "deny", "cancel", "failure", "expire", "refund", "partial_refund":
		return true
	}
	return false
}
