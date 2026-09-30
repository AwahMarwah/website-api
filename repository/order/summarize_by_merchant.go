package order

import "gorm.io/gorm"

// summarizeByMerchantSQL menghitung jumlah order per status beserta pendapatan dari
// item milik seller tersebut.
//
// COUNT(DISTINCT o.id) dipakai karena satu order bisa memuat banyak item dari seller
// yang sama, dan tanpa DISTINCT order akan terhitung berkali-kali.
const summarizeByMerchantSQL = `
SELECT o.status                                             AS status,
       COUNT(DISTINCT o.id)                                AS total,
       COALESCE(SUM(oi.subtotal), 0)                       AS revenue
FROM orders o
JOIN order_items oi ON oi.order_id = o.id AND oi.merchant_id = ?
GROUP BY o.status
`

// SummarizeByMerchant mengembalikan agregat order per status untuk satu merchant.
func (r *repo) SummarizeByMerchant(merchantID string) *gorm.DB {
	return r.db.Raw(summarizeByMerchantSQL, merchantID)
}
