package order

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"website-api/model/order"
	midtransProvider "website-api/third-party/provider/midtrans"

	"github.com/midtrans/midtrans-go"
	"gorm.io/gorm"
)

// CreatePaymentLink membuat payment link Midtrans untuk order yang masih PENDING.
func (s *service) CreatePaymentLink(orderID, userID, roleName string) (order.PaymentLinkResponse, int, error) {
	var resData order.PaymentLinkResponse

	existingOrder, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resData, http.StatusNotFound, fmt.Errorf("order not found")
		}
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal mengambil order: %w", err)
	}

	// Ownership check (IDOR) - hanya pemilik order atau super_admin/admin
	if existingOrder.UserID != userID && roleName != "super_admin" && roleName != "admin" {
		return resData, http.StatusForbidden, fmt.Errorf("forbidden")
	}

	if existingOrder.Status != "PENDING" {
		return resData, http.StatusConflict, fmt.Errorf("order sudah tidak dalam status pending")
	}

	items, err := s.orderRepo.FindItemsByOrderID(orderID)
	if err != nil {
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal mengambil item order: %w", err)
	}

	var itemDetails []midtransProvider.PaymentLinkItemDetail
	for _, item := range items {
		itemDetails = append(itemDetails, midtransProvider.PaymentLinkItemDetail{
			ID:       item.ID,
			Name:     item.ProductVariantID,
			Price:    int64(item.Price),
			Quantity: int32(item.Qty),
		})
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	reqLink := midtransProvider.PaymentLinkRequest{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  existingOrder.ID,
			GrossAmt: int64(existingOrder.TotalAmount),
		},
		Items: itemDetails,
		Expiry: &midtransProvider.ExpiryUnit{
			StartTime: time.Now().Format("2006-01-02 15:04:05 -0700"),
			Unit:      "hours",
			Duration:  24,
		},
	}

	linkResp, err := s.midtransProvider.CreatePaymentLink(reqLink)
	if err != nil {
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal membuat payment link: %w", err)
	}

	updateMap := map[string]interface{}{
		"payment_url": linkResp.URL,
		"expired_at":  expiresAt,
	}
	if err := s.orderRepo.UpdatePaymentInfo(orderID, updateMap); err != nil {
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal menyimpan payment link: %w", err)
	}

	resData = order.PaymentLinkResponse{
		OrderID: existingOrder.ID,
		URL:     linkResp.URL,
	}

	return resData, http.StatusOK, nil
}