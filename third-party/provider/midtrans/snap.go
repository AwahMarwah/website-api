package midtrans

import (
	"github.com/midtrans/midtrans-go/snap"
)

// CreateTransaction membuat Snap transaction dan mengembalikan token + redirect url
func (c *Client) CreateTransaction(req *snap.Request) (*snap.Response, error) {
	resp, mErr := c.snapClient.CreateTransaction(req)
	if mErr != nil {
		return resp, ConvertMidtransError(mErr)
	}
	return resp, nil
}
