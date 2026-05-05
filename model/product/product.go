package product

import "time"

type Product struct {
	Id          string
	BrandId     string
	Sku         string
	Name        string
	Slug        string
	Description string
	BasePrice   float32
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
