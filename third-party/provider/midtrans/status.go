package midtrans

import (
	"github.com/midtrans/midtrans-go/coreapi"
)

// GetTransactionStatus mengambil status transaksi dari Midtrans berdasarkan order id
func (c *Client) GetTransactionStatus(orderID string) (*coreapi.TransactionStatusResponse, error) {
	resp, mErr := c.coreapiClient.CheckTransaction(orderID)
	if mErr != nil {
		return resp, ConvertMidtransError(mErr)
	}
	return resp, nil
}
