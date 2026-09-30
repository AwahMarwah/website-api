package midtrans

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// ProviderName identifies the payment provider used for refund records.
const ProviderName = "midtrans"

// refundPath template untuk endpoint refund per-status transaksi.
// https://docs.midtrans.com/reference/refund-transaction
var refundPath = map[string]string{
	"pending":    "pending",
	"settlement": "settlement",
	"capture":    "capture",
}

// Refund mengembalikan dana ke pembeli untuk order tertentu.
//
// amount yang dikirim adalah nominal refund, bukan sisa refundable. Midtrans
// menjumlahkan refund sebelumnya untuk order yang sama, jadi Aman mengirim
// nominal parsial tanpa menghitung total yang sudah dikembalikan.
func (c *Client) Refund(orderID, transactionStatus, reason string, amount int64) (*RefundResponse, error) {
	path, ok := refundPath[transactionStatus]
	if !ok {
		return nil, fmt.Errorf("transaction status %q tidak mendukung refund langsung", transactionStatus)
	}
	if c.serverKey == "" {
		return nil, fmt.Errorf("server key is not configured")
	}

	jsonReq, err := json.Marshal(RefundRequest{
		TransactionStatus: transactionStatus,
		RefundAmount:     amount,
		Reason:           reason,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal refund request: %w", err)
	}

	key := c.serverKey
	resp := &RefundResponse{}

	httpErr := c.midtransClient.Call(
		http.MethodPost,
		fmt.Sprintf("%s/v2/%s/%s/refund", c.env.BaseUrl(), path, orderID),
		&key,
		c.options,
		bytes.NewBuffer(jsonReq),
		resp,
	)
	if httpErr != nil {
		return nil, ConvertMidtransError(httpErr)
	}
	return resp, nil
}
