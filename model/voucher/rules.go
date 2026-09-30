package voucher

import (
	"errors"
	"fmt"
	"math"
	"time"

	"website-api/common"
)

// Aturan voucher diletakkan di model, bukan di layer service, supaya pratinjau di
// admin dan checkout tidak pernah menghitung diskon dengan cara berbeda.

// Validate memeriksa syarat voucher yang bisa dinilai tanpa database: status aktif,
// periode berlaku, dan kecukupan nilai belanja.
//
// Syarat yang butuh database (kuota terpakai, jumlah pemakaian per pengguna) tetap
// dicek di layer service karena membutuhkan repo.
func (v Voucher) Validate(now time.Time, amount float64) error {
	if !v.IsActive {
		return NewValidationError("voucher tidak aktif")
	}
	if v.StartsAt != nil && now.Before(*v.StartsAt) {
		return NewValidationError("voucher belum mulai berlaku")
	}
	if v.EndsAt != nil && now.After(*v.EndsAt) {
		return NewValidationError("voucher sudah berakhir")
	}
	if v.Quota != nil && v.UsedCount >= *v.Quota {
		return NewValidationError("kuota voucher sudah habis")
	}
	if amount < v.MinSpend {
		return NewValidationError("minimum belanja untuk voucher ini adalah %.0f", v.MinSpend)
	}
	if v.ComputeDiscount(amount) <= 0 {
		return NewValidationError("voucher tidak berlaku untuk nilai belanja ini")
	}
	return nil
}

// ComputeDiscount menghitung nominal diskon terhadap nilai belanja.
// PERCENT dibaca sebagai persentase, FIXED sebagai nominal rupiah langsung.
// Diskon tidak pernah melebihi nilai belanja supaya total tidak bisa negatif.
func (v Voucher) ComputeDiscount(amount float64) float64 {
	var discount float64

	if v.Type == common.DiscountTypePercent {
		discount = amount * v.Value / 100
		if v.MaxDiscount != nil && discount > *v.MaxDiscount {
			discount = *v.MaxDiscount
		}
	} else {
		discount = v.Value
	}

	if discount > amount {
		discount = amount
	}
	if discount < 0 {
		discount = 0
	}
	return math.Round(discount*100) / 100
}

// ValidationError adalah kegagalan validasi yang pesannya aman ditampilkan
// langsung ke pengguna, berbeda dari error internal.
type ValidationError struct {
	Message string
}

func NewValidationError(format string, args ...any) *ValidationError {
	return &ValidationError{Message: fmt.Sprintf(format, args...)}
}

func (e *ValidationError) Error() string { return e.Message }

// AsValidationError melacak apakah error berasal dari validasi voucher,
// supaya service bisa membedakan pesan untuk pengguna dan kegagalan internal.
func AsValidationError(err error) (*ValidationError, bool) {
	var verr *ValidationError
	if errors.As(err, &verr) {
		return verr, true
	}
	return nil, false
}
