package refund

import (
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"website-api/library/pagination"
	refundModel "website-api/model/refund"
	midtransProvider "website-api/third-party/provider/midtrans"
)

// refundResponse adalah alias lokal supaya berkas ini tidak perlu mengimpor
// package provider hanya untuk satu tipe.
type refundResponse = midtransProvider.RefundResponse

func (s *service) List(reqQuery *refundModel.ListRefundReqQuery) ([]refundModel.RefundResponse, int64, int, error) {
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)

	list, count, err := s.refundRepo.List(reqQuery)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, fmt.Errorf("gagal mengambil daftar refund: %w", err)
	}

	res := make([]refundModel.RefundResponse, 0, len(list))
	for _, v := range list {
		res = append(res, mapToResponse(v))
	}
	return res, count, http.StatusOK, nil
}

func (s *service) Detail(refundID string) (refundModel.RefundResponse, int, error) {
	v, err := s.refundRepo.FindByID(refundID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return refundModel.RefundResponse{}, http.StatusNotFound, fmt.Errorf("refund not found")
		}
		return refundModel.RefundResponse{}, http.StatusInternalServerError,
			fmt.Errorf("gagal mengambil refund: %w", err)
	}
	return mapToResponse(v), http.StatusOK, nil
}

func mapToResponse(v refundModel.Refund) refundModel.RefundResponse {
	return refundModel.RefundResponse{
		ID:               v.ID,
		OrderID:          v.OrderID,
		MerchantID:       v.MerchantID,
		RequestedBy:      v.RequestedBy,
		Amount:           v.Amount,
		Reason:           v.Reason,
		Status:           v.Status,
		Provider:         v.Provider,
		ProviderRefundID: v.ProviderRefundID,
		ApprovedBy:       v.ApprovedBy,
		Restocked:        v.Restocked,
		FailureReason:    v.FailureReason,
		CreatedAt:        v.CreatedAt,
		UpdatedAt:        v.UpdatedAt,
	}
}
