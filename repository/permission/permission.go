package permission

import (
	permissionModel "website-api/model/permission"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		FindNamesByRole(roleID string) ([]string, error)
		FindAll() ([]permissionModel.Permission, error)
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}

// FindNamesByRole mengembalikan daftar nama permission yang dimiliki role
// melalui relasi role_menus -> menus -> permissions.menu_id.
func (r *repo) FindNamesByRole(roleID string) ([]string, error) {
	var names []string
	err := r.db.Model(&permissionModel.Permission{}).
		Select("permissions.name").
		Joins("JOIN menus m ON m.id = permissions.menu_id").
		Joins("JOIN role_menus rm ON rm.menu_id = m.id").
		Where("rm.role_id = ? AND m.deleted_at IS NULL", roleID).
		Pluck("permissions.name", &names).Error
	return names, err
}

func (r *repo) FindAll() ([]permissionModel.Permission, error) {
	var permissions []permissionModel.Permission
	err := r.db.Find(&permissions).Error
	return permissions, err
}
