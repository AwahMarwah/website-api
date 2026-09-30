package refund

import (
	"database/sql"

	"gorm.io/gorm"

	refundModel "website-api/model/refund"
)

type (
	IRepo interface {
		Create(refundModel.Refund) error
		FindByID(id string) (refundModel.Refund, error)
		FindByOrderID(orderID string) ([]refundModel.Refund, error)
		FindPendingByOrderID(orderID string) ([]refundModel.Refund, error)
		List(reqQuery *refundModel.ListRefundReqQuery) ([]refundModel.Refund, int64, error)
		UpdateStatus(id, status string, values map[string]any) error
		SumCompletedByOrder(orderID string) (float64, error)
		WithTx(tx *gorm.DB) IRepo
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}

func (r *repo) WithTx(tx *gorm.DB) IRepo {
	if tx == nil {
		return r
	}
	return &repo{db: tx}
}

func (r *repo) Create(v refundModel.Refund) error {
	return r.db.Create(&v).Error
}

func (r *repo) FindByID(id string) (refundModel.Refund, error) {
	var v refundModel.Refund
	err := r.db.Where("id = ?", id).First(&v).Error
	if err == sql.ErrNoRows {
		return v, gorm.ErrRecordNotFound
	}
	return v, err
}

func (r *repo) FindByOrderID(orderID string) ([]refundModel.Refund, error) {
	var list []refundModel.Refund
	return list, r.db.Where("order_id = ?", orderID).Order("created_at ASC").Find(&list).Error
}

func (r *repo) FindPendingByOrderID(orderID string) ([]refundModel.Refund, error) {
	var list []refundModel.Refund
	err := r.db.Where("order_id = ? AND status IN ?", orderID,
		[]string{refundModel.StatusPending, refundModel.StatusApproved, refundModel.StatusProcessing}).
		Find(&list).Error
	return list, err
}

func (r *repo) List(reqQuery *refundModel.ListRefundReqQuery) (list []refundModel.Refund, count int64, err error) {
	query := r.db.Model(&refundModel.Refund{})

	if reqQuery.Status != "" {
		query = query.Where("status = ?", reqQuery.Status)
	}
	if reqQuery.OrderID != "" {
		query = query.Where("order_id = ?", reqQuery.OrderID)
	}
	if reqQuery.MerchantID != "" {
		query = query.Where("merchant_id = ?", reqQuery.MerchantID)
	}

	if err = query.Count(&count).Error; err != nil {
		return nil, count, err
	}
	err = query.Limit(reqQuery.Limit).Offset(reqQuery.Offset).
		Order("created_at DESC").Find(&list).Error
	return list, count, err
}

func (r *repo) UpdateStatus(id, status string, values map[string]any) error {
	values["status"] = status
	return r.db.Model(&refundModel.Refund{}).Where("id = ?", id).Updates(values).Error
}

// SumCompletedByOrder menjumlahkan refund yang sudah berhasil untuk satu order,
// supaya total pengembalian tidak bisa melebihi nilai order.
func (r *repo) SumCompletedByOrder(orderID string) (float64, error) {
	var total float64
	err := r.db.Model(&refundModel.Refund{}).
		Where("order_id = ? AND status = ?", orderID, refundModel.StatusCompleted).
		Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}
