package product

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductImage struct {
	ID        string
	ProductID string
	ImageURL  string
	IsPrimary bool
	SortOrder int
	CreatedAt time.Time
}

type ImageSortOrder struct {
	ID        string
	SortOrder int
}

func (img *ProductImage) BeforeCreate(*gorm.DB) error {
	img.ID = uuid.New().String()
	return nil
}