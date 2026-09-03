package menu

import (
	menu "website-api/model/menu"

	"gorm.io/gorm"
)

func (r *repo) Create(m *menu.Menu) error {
	return r.db.Create(m).Error
}

func (r *repo) Update(m *menu.Menu) error {
	return r.db.Save(m).Error
}

func (r *repo) SoftDelete(id string) error {
	return r.db.Model(&menu.Menu{}).Where("id = ?", id).Updates(map[string]interface{}{
		"deleted_at": gorm.Expr("NOW()"),
	}).Error
}

func (r *repo) FindByID(id string) (menu.Menu, error) {
	var m menu.Menu
	err := r.db.Where("id = ?", id).First(&m).Error
	return m, err
}

func (r *repo) FindAll() ([]menu.Menu, error) {
	var menus []menu.Menu
	err := r.db.Where("deleted_at IS NULL").Order("sort_order ASC").Find(&menus).Error
	return menus, err
}

func (r *repo) FindWithPagination(limit, offset int) ([]menu.Menu, int64, error) {
	var menus []menu.Menu
	var count int64
	if err := r.db.Model(&menu.Menu{}).Where("deleted_at IS NULL").Count(&count).Error; err != nil {
		return nil, count, err
	}
	err := r.db.Where("deleted_at IS NULL").Order("sort_order ASC").Limit(limit).Offset(offset).Find(&menus).Error
	return menus, count, err
}

func (r *repo) FindMenusByRole(roleID string) ([]menu.Menu, error) {
	var menus []menu.Menu
	err := r.db.
		Joins("JOIN role_menus rm ON rm.menu_id = menus.id").
		Where("rm.role_id = ? AND menus.deleted_at IS NULL", roleID).
		Order("menus.sort_order ASC").
		Find(&menus).Error
	return menus, err
}
