package merchant

import (
	"fmt"

	"website-api/common"
)

// summarizeOrders menghitung jumlah order per status dan total pendapatan dari item
// milik seller, dalam satu query.
//
// Pendapatan dihitung dari order_items milik seller, bukan dari orders.total_amount.
// Total order ikut mencakup ongkir seller lain dan diskon, jadi bukan angka yang
// boleh dipakai untuk membayar seller.
func (s *service) summarizeOrders(merchantID string) (map[string]int64, float64, error) {
	type orderStatRow struct {
		Status  string
		Total   int64
		Revenue float64
	}

	var rows []orderStatRow
	if err := s.orderRepo.SummarizeByMerchant(merchantID).Scan(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("gagal merangkum order: %w", err)
	}

	counts := make(map[string]int64, len(rows))
	var revenue float64
	for _, row := range rows {
		counts[row.Status] = row.Total
		switch row.Status {
		case common.OrderStatusPaid,
			common.OrderStatusProcessing,
			common.OrderStatusShipped,
			common.OrderStatusCompleted:
			// Hanya order yang sudah dibayar yang dihitung sebagai pendapatan:
			// order PENDING bisa dibatalkan dan stoknya dikembalikan.
			revenue += row.Revenue
		}
	}
	return counts, revenue, nil
}
