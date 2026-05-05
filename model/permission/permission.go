package permission

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID          string
	Name        string
	DisplayName string
	Description string
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   *time.Time
	UpdatedBy   *string
	DeletedAt   *time.Time
	DeletedBy   *string
}

func (permission *Permission) BeforeCreate(*gorm.DB) error {
	permission.ID = uuid.New().String()
	return nil
}
