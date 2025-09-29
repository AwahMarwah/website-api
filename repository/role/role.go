package role

import (
	"gorm.io/gorm"
	modelRole "website-api/model/role"
)

type (
	IRepo interface {
		Create(reqBody *modelRole.Role) error
		Find() (roles []modelRole.Role, err error)
		Take(selectParams []string, condition *modelRole.Role) (role modelRole.Role, err error)
		Update(id *string, values *map[string]any) (err error)
		Seed() error
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}
