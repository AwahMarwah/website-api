package product_variant

// findMerchantIDsByVariantIDsSQL memetakan variant ke merchant pemiliknya dalam
// satu query. Dipakai saat memvalidasi voucher milik seller: keranjang harus
// benar-benar memuat barang dari seller tersebut.
const findMerchantIDsByVariantIDsSQL = `
SELECT DISTINCT COALESCE(p.merchant_id, '') AS merchant_id
FROM product_variants pv
JOIN products p ON p.id = pv.product_id
WHERE pv.id IN ?
`

func (r *repo) FindMerchantIDsByVariantIDs(variantIDs []string) ([]string, error) {
	if len(variantIDs) == 0 {
		return []string{}, nil
	}
	var merchantIDs []string
	err := r.db.Raw(findMerchantIDsByVariantIDsSQL, variantIDs).Scan(&merchantIDs).Error
	return merchantIDs, err
}
