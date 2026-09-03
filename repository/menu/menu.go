package menu

import (
	"time"
	menu "website-api/model/menu"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		Create(m *menu.Menu) error
		Update(m *menu.Menu) error
		SoftDelete(id string) error
		FindByID(id string) (menu.Menu, error)
		FindAll() ([]menu.Menu, error)
		FindWithPagination(limit, offset int) ([]menu.Menu, int64, error)
		AssignMenus(roleID string, menuIDs []string, now time.Time) error
		FindMenusByRole(roleID string) ([]menu.Menu, error)
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
