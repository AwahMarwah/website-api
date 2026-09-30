package order

import (
	"time"

	"website-api/model/merchant"
	"website-api/model/product-variant"
)

type Order struct {
	ID          string
	UserID      string
	AddressID   string
	TotalAmount float64
	ShippingFee float64
	// DiscountAmount sudah dipotong di dalam TotalAmount; disimpan terpisah agar
	// voucher yang berakhir tidak mengubah tampilan order lama.
	DiscountAmount float64
	Status         string
	PaymentMethod  string
	PaymentToken   string
	PaymentURL     *string
	ExpiredAt      *time.Time
	// Snapshot penerima pada saat checkout. AddressID tetap disimpan sebagai referensi,
	// tapi alamat bisa diedit/dihapus setelah itu sehingga tidak boleh jadi sumber tampilan.
	AddressSnapshot AddressSnapshot
	BuyerName       string
	BuyerEmail      string
	BuyerPhone      string
	Note            *string
	PaidAt          *time.Time
	ShippedAt       *time.Time
	CompletedAt     *time.Time
	CancelledAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Items           []OrderItem     `gorm:"foreignKey:OrderID"`
	StatusHistories []OrderStatusHistory
}

// AddressSnapshot adalah salinan user_addresses yang dibekukan saat checkout.
// Hanya memuat kolom yang benar-benar ada di user_addresses; nama wilayah
// (provinsi/kecamatan) tidak disimpan karena tabel itu hanya menyimpan ID.
type AddressSnapshot struct {
	RecipientName string `json:"recipient_name"`
	PhoneNumber   string `json:"phone_number"`
	FullAddress   string `json:"full_address"`
	City          string `json:"city"`
	PostalCode    string `json:"postal_code"`
	ProvinceID    string `json:"province_id"`
	CityID        string `json:"city_id"`
	DistrictID    string `json:"district_id"`
	SubdistrictID string `json:"subdistrict_id"`
	DestinationID int64  `json:"destination_id"`
}

type OrderItem struct {
	ID               string
	OrderID          string
	ProductVariantID string
	// MerchantID adalah pemilik barang saat transaksi dibuat. Inilah yang membuat
	// seller bisa melihat order-nya tanpa join ke tabel produk.
	MerchantID string
	Price      float64
	Qty        int
	Subtotal   float64
	// Snapshot tampilan: produk bisa berubah atau di-soft delete setelah order dibuat.
	ProductName     string
	VariantName     string
	ProductImageURL string
	Sku             string
	TotalWeightGram int
	ProductVariant  *product_variant.ProductVariant `gorm:"foreignKey:ProductVariantID"`
	Merchant        *merchant.Merchant              `gorm:"foreignKey:MerchantID"`
}

type OrderStatusHistory struct {
	ID         uint
	OrderID    string
	FromStatus *string
	ToStatus   string
	ActorID    *string
	ActorRole  *string
	Note       *string
	CreatedAt  time.Time
}

type OrderMerchantShipping struct {
	ID         string
	OrderID    string
	MerchantID string
	Courier    string
	Service    string
	Cost       int64
	Etd        string
	WeightGram int
	CreatedAt  time.Time
	Merchant   *merchant.Merchant `gorm:"foreignKey:MerchantID"`
}
