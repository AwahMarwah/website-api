package cart

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItem struct {
	ID               string
	UserID           string
	ProductVariantID string
	Qty              int
	CreatedAt        time.Time
}

func (cartItem *CartItem) BeforeCreate(*gorm.DB) error {
	cartItem.ID = uuid.New().String()
	return nil
}
