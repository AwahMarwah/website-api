package voucher

import (
	"website-api/database/transaction"
	voucherModel "website-api/model/voucher"
	"website-api/repository/voucher"
)

type IService interface {
	Create(req *voucherModel.CreateVoucherReq) (id string, statusCode int, err error)
	List(reqQuery *voucherModel.ListVoucherReqQuery) (vouchers []voucherModel.VoucherResponse, count int64, statusCode int, err error)
	Detail(id string) (voucherModel.VoucherResponse, int, error)
	Update(id string, req *voucherModel.UpdateVoucherReq) (int, error)
	Delete(id string) (int, error)
	// Validate mengecek voucher terhadap nilai belanja dan mengembalikan diskonnya,
	// dipakai untuk pratinjau di frontend sebelum checkout.
	Validate(code, userID string, amount float64) (voucherModel.ValidateVoucherResponse, error)
}

type service struct {
	voucherRepo voucher.IRepo
	txManager   transaction.ITransactionManager
}

func NewService(voucherRepo voucher.IRepo, txManager transaction.ITransactionManager) IService {
	return &service{voucherRepo: voucherRepo, txManager: txManager}
}
