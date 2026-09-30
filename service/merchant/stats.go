package merchant

import (
	"fmt"
	"net/http"

	"website-api/common"
	merchantModel "website-api/model/merchant"
)

// Stats menghitung ringkasan dashboard seller: jumlah produk per status, jumlah order
// per status, pendapatan dari item miliknya, dan saldo setelah komisi.
func (s *service) Stats(userID string) (merchantModel.SellerStatsResponse, int, error) {
	merchantID, statusCode, err := s.resolveMerchantID(userID)
	if err != nil {
		return merchantModel.SellerStatsResponse{}, statusCode, err
	}

	var res merchantModel.SellerStatsResponse

	// Satu query agregat untuk ketiga status produk, bukan tiga query terpisah.
	var productCounts []struct {
		Status string
		Total  int64
	}
	if err := s.productRepo.CountByStatus(merchantID).Scan(&productCounts).Error; err != nil {
		return res, http.StatusInternalServerError, fmt.Errorf("gagal menghitung produk: %w", err)
	}
	for _, row := range productCounts {
		res.TotalProducts += row.Total
		switch row.Status {
		case common.ProductStatusActive:
			res.ActiveProducts = row.Total
		case common.ProductStatusDraft:
			res.DraftProducts = row.Total
		}
	}

	orderCounts, revenue, err := s.summarizeOrders(merchantID)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}
	res.TotalOrders = orderCounts[common.OrderStatusPending] +
		orderCounts[common.OrderStatusPaid] +
		orderCounts[common.OrderStatusProcessing] +
		orderCounts[common.OrderStatusShipped] +
		orderCounts[common.OrderStatusCompleted] +
		orderCounts[common.OrderStatusCancelled] +
		orderCounts[common.OrderStatusRefunded]
	res.PendingOrders = orderCounts[common.OrderStatusPending]
	res.ProcessingOrders = orderCounts[common.OrderStatusProcessing]
	res.ShippedOrders = orderCounts[common.OrderStatusShipped]
	res.CompletedOrders = orderCounts[common.OrderStatusCompleted]
	res.ItemRevenue = revenue

	balance, err := s.settlementRepo.GetBalance(merchantID)
	if err != nil {
		return res, http.StatusInternalServerError, fmt.Errorf("gagal mengambil saldo: %w", err)
	}
	res.Balance = balance

	// CommissionEarned positif untuk ditampilkan, sementara di ledger disimpan negatif.
	commission, err := s.settlementRepo.SumCommission(merchantID)
	if err != nil {
		return res, http.StatusInternalServerError, fmt.Errorf("gagal menghitung komisi: %w", err)
	}
	res.CommissionEarned = -commission

	return res, http.StatusOK, nil
}
