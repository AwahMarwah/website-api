package midtrans

import (
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type Provider interface {
	// CreateTransaction membuat Snap transaction dan mengembalikan token + redirect url
	CreateTransaction(req *snap.Request) (*snap.Response, error)
	// GetTransactionStatus mengambil status transaksi dari Midtrans
	GetTransactionStatus(orderID string) (*coreapi.TransactionStatusResponse, error)
	// CreatePaymentLink membuat payment link untuk dibagikan
	CreatePaymentLink(req PaymentLinkRequest) (*PaymentLinkResponse, error)
}

// ProviderLogger menyediakan akses helper yang berkaitan dengan provider
type ProviderLogger interface {
	ServerKey() string
	IsProduction() bool
}
