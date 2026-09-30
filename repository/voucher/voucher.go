package voucher

import (
	"gorm.io/gorm"

	voucherModel "website-api/model/voucher"
)

type (
	IRepo interface {
		FindByID(id string) (voucherModel.Voucher, error)
		FindByCode(code string) (voucherModel.Voucher, error)
		List(reqQuery *voucherModel.ListVoucherReqQuery) (vouchers []voucherModel.Voucher, count int64, err error)
		Create(voucherModel.Voucher) error
		Update(id string, values map[string]any) error
		Delete(id string) error
		// ConsumeQuota menambah used_count hanya selama kuota masih tersedia.
		// RowsAffected == 0 berarti kuota habis, dipakai untuk menutup race saat checkout bersamaan.
		ConsumeQuota(id string) (int64, error)
		ReleaseQuota(id string) error
		CountUserRedemptions(voucherID, userID string) (int64, error)
		CreateRedemption(redemption voucherModel.VoucherRedemption) error
		RebindProducts(voucherID string, productIDs []string) error
		FindProductIDs(voucherID string) ([]string, error)
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
