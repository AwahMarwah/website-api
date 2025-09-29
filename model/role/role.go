package role

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type (
	Role struct {
		Id          string
		Name        string
		DisplayName string
		Description string
		IsActive    bool
		CreatedAt   time.Time
		CreatedBy   string
		UpdatedAt   *time.Time
		UpdatedBy   string
		DeletedAt   *time.Time
		DeletedBy   string
	}
)

func (role *Role) BeforeCreate(*gorm.DB) error {
	role.Id = uuid.New().String()
	return nil
}
