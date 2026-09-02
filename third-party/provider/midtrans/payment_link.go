package midtrans

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreatePaymentLink membuat payment link Midtrans untuk dibagikan ke pelanggan.
// https://docs.midtrans.com/reference/create-payment-link
func (c *Client) CreatePaymentLink(req PaymentLinkRequest) (*PaymentLinkResponse, error) {
	if c.paymentLinkKey == "" && c.serverKey == "" {
		return nil, fmt.Errorf("payment link key / server key is not configured")
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payment link request: %w", err)
	}

	key := c.paymentLinkKeyOrDefault()
	resp := &PaymentLinkResponse{}

	httpErr := c.midtransClient.Call(
		http.MethodPost,
		fmt.Sprintf("%s/v1/payment-links", c.env.BaseUrl()),
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
