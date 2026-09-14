package product

import productModel "website-api/model/product"

func (r *repo) CreateProductImage(img *productModel.ProductImage) error {
	return r.db.Create(img).Error
}

func (r *repo) ProductImageCount(productID string) (int64, error) {
	var count int64
	err := r.db.Model(&productModel.ProductImage{}).Where("product_id = ?", productID).Count(&count).Error
	return count, err
}

func (r *repo) FindImagesByProductID(productID string) ([]productModel.ProductImage, error) {
	var images []productModel.ProductImage
	err := r.db.Where("product_id = ?", productID).Order("is_primary DESC, sort_order ASC").Find(&images).Error
	return images, err
}

func (r *repo) MaxSortOrder(productID string) (int, error) {
	var max int
	err := r.db.Model(&productModel.ProductImage{}).Where("product_id = ?", productID).Select("COALESCE(MAX(sort_order),0)").Scan(&max).Error
	return max, err
}

func (r *repo) SetPrimary(productID, imageID string) error {
	// demote semua gambar produk ke non-primary, lalu set gambar tsb primary
	if err := r.db.Model(&productModel.ProductImage{}).Where("product_id = ?", productID).Update("is_primary", false).Error; err != nil {
		return err
	}
	return r.db.Model(&productModel.ProductImage{}).Where("id = ?", imageID).Update("is_primary", true).Error
}

func (r *repo) DeleteImage(imageID string) error {
	return r.db.Where("id = ?", imageID).Delete(&productModel.ProductImage{}).Error
}

func (r *repo) SetImagePrimaryFalse(imageID string) error {
	return r.db.Model(&productModel.ProductImage{}).Where("id = ?", imageID).Update("is_primary", false).Error
}

func (r *repo) UpdateSortOrders(productID string, orders []productModel.ImageSortOrder) error {
	for _, o := range orders {
		if err := r.db.Model(&productModel.ProductImage{}).
			Where("id = ? AND product_id = ?", o.ID, productID).
			Update("sort_order", o.SortOrder).Error; err != nil {
			return err
		}
	}
	return nil
}