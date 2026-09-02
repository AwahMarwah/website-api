package midtrans

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type Client struct {
	snapClient      snap.Client
	coreapiClient   coreapi.Client
	midtransClient  midtrans.HttpClient
	serverKey       string
	clientKey       string
	paymentLinkKey  string
	env             midtrans.EnvironmentType
	options         *midtrans.ConfigOptions
	isProduction    bool
}

func NewClient() *Client {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	clientKey := os.Getenv("MIDTRANS_CLIENT_KEY")
	paymentLinkKey := os.Getenv("MIDTRANS_PAYMENT_LINK_KEY")
	isProduction := strings.TrimSpace(os.Getenv("MIDTRANS_IS_PRODUCTION")) == "true"

	env := midtrans.Sandbox
	if isProduction {
		env = midtrans.Production
	}

	snapClient := snap.Client{}
	snapClient.New(serverKey, env)

	coreClient := coreapi.Client{}
	coreClient.New(serverKey, env)
	coreClient.ClientKey = clientKey

	return &Client{
		snapClient:     snapClient,
		coreapiClient:  coreClient,
		midtransClient: midtrans.GetHttpClient(env),
		serverKey:      serverKey,
		clientKey:      clientKey,
		paymentLinkKey: paymentLinkKey,
		env:            env,
		options: &midtrans.ConfigOptions{
			PaymentOverrideNotification: midtrans.PaymentOverrideNotification,
			PaymentAppendNotification:   midtrans.PaymentAppendNotification,
		},
		isProduction: isProduction,
	}
}

func (c *Client) ServerKey() string {
	return c.serverKey
}

func (c *Client) ClientKey() string {
	return c.clientKey
}

func (c *Client) IsProduction() bool {
	return c.isProduction
}

// ConvertMidtransError mengubah *midtrans.Error menjadi error standar
func ConvertMidtransError(mErr *midtrans.Error) error {
	if mErr == nil {
		return nil
	}
	if mErr.RawError != nil {
		return mErr.RawError
	}
	return mErr
}

// errorFromRaw mengubah raw body dari request HTTP menjadi error
func errorFromRaw(status int, body []byte) error {
	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err == nil {
		if msg, ok := parsed["error_message"]; ok {
			return fmt.Errorf("midtrans returned status %d: %v", status, msg)
		}
	}
	return fmt.Errorf("midtrans returned status %d: %s", status, string(body))
}

// paymentLinkKeyOrDefault mengembalikan payment link key atau server key sebagai fallback
func (c *Client) paymentLinkKeyOrDefault() string {
	if c.paymentLinkKey != "" {
		return c.paymentLinkKey
	}
	return c.serverKey
}
