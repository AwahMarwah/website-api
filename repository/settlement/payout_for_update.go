package settlement

import (
	"database/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	settlementModel "website-api/model/settlement"
)

// FindPayoutByIDForUpdate mengambil payout dengan SELECT ... FOR UPDATE.
// Dipakai saat approval supaya dua approver yang bersamaan tidak menghasilkan
// dua kali pencairan.
func (r *repo) FindPayoutByIDForUpdate(id string) (settlementModel.Payout, error) {
	var payout settlementModel.Payout
	err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&payout).Error
	if err == sql.ErrNoRows {
		return payout, gorm.ErrRecordNotFound
	}
	return payout, err
}
