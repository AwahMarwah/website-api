package settlement

import (
	"fmt"
	"net/http"

	"website-api/library/pagination"
	settlementModel "website-api/model/settlement"
)

func (s *service) ListLedgers(merchantID string, reqQuery *settlementModel.ListLedgerReqQuery) ([]settlementModel.MerchantLedgerResponse, int64, int, error) {
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)

	ledgers, count, err := s.settlementRepo.ListLedgers(merchantID, reqQuery.Type, reqQuery.Limit, reqQuery.Offset)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, fmt.Errorf("gagal mengambil ledger: %w", err)
	}

	res := make([]settlementModel.MerchantLedgerResponse, 0, len(ledgers))
	for _, l := range ledgers {
		res = append(res, settlementModel.MerchantLedgerResponse{
			ID:           l.ID,
			MerchantID:   l.MerchantID,
			OrderID:      l.OrderID,
			Type:         l.Type,
			Amount:       l.Amount,
			BalanceAfter: l.BalanceAfter,
			Note:         l.Note,
			CreatedAt:    l.CreatedAt,
		})
	}
	return res, count, http.StatusOK, nil
}

func (s *service) Balance(merchantID string) (settlementModel.BalanceResponse, int, error) {
	balance, err := s.settlementRepo.GetBalance(merchantID)
	if err != nil {
		return settlementModel.BalanceResponse{}, http.StatusInternalServerError,
			fmt.Errorf("gagal mengambil saldo: %w", err)
	}

	pending, err := s.settlementRepo.SumPayout(merchantID)
	if err != nil {
		return settlementModel.BalanceResponse{}, http.StatusInternalServerError,
			fmt.Errorf("gagal mengambil total payout: %w", err)
	}

	return settlementModel.BalanceResponse{
		MerchantID: merchantID,
		Balance:    balance,
		PendingPayout: pending,
	}, http.StatusOK, nil
}
