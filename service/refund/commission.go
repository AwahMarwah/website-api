package refund

import (
	"gorm.io/gorm"

	refundModel "website-api/model/refund"
)

// commissionReversed menghitung berapa komisi yang perlu dikembalikan ke seller untuk
// satu refund.
//
// 수수ari dihitung secara proporsional terhadap nominal refund, bukan dari seluruh
// order, supaya refund parsial hanya mengembalikan komisi sebesar bagian yang benar-benar
// dikembalikan ke pembeli.
func (s *service) commissionReversed(tx *gorm.DB, row refundModel.Refund) (float64, error) {
	if row.MerchantID == nil || *row.MerchantID == "" {
		return 0, nil
	}

	repo := s.settlementRepo.WithTx(tx)
	orderRow, err := s.orderRepo.FindByID(row.OrderID)
	if err != nil {
		return 0, err
	}
	if orderRow.TotalAmount <= 0 {
		return 0, nil
	}

	// SumCommission mengembalikan nilai positif sementara di ledger tersimpan negatif.
	commission, err := repo.SumCommission(*row.MerchantID)
	if err != nil {
		return 0, err
	}

	ratio := row.Amount / orderRow.TotalAmount
	return commission * ratio, nil
}
