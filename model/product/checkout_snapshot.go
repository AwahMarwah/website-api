package product

// CheckoutSnapshot adalah data minimum produk yang dibutuhkan saat checkout:
// nama, pemilik barang, dan gambar utama. Satu query, bukan tiga.
//
// Sengaja tidak memakai struct Product penuh: checkout cukup sering dipanggil dan
// tidak butuh seluruh relasi produk.
type CheckoutSnapshot struct {
	ProductID   string
	ProductName string
	MerchantID  string
	ImageURL    string
}
