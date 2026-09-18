package product

import productModel "website-api/model/product"

func (r *repo) Create(product productModel.Product) error {
	return r.db.Create(&product).Error
}

func (r *repo) UpdateProduct(id string, values map[string]any) error {
	return r.db.Model(&productModel.Product{}).Where("id = ?", id).Updates(values).Error
}

func (r *repo) SetInactive(id string) error {
	return r.db.Model(&productModel.Product{}).Where("id = ?", id).Update("status", "inactive").Error
}

func (r *repo) RebindCategories(productID string, catIDs []string) error {
	// hapus relasi lama lalu insert baru
	if err := r.db.Exec("DELETE FROM product_categories WHERE product_id = ?", productID).Error; err != nil {
		return err
	}
	for _, cid := range catIDs {
		if err := r.db.Exec(
			"INSERT INTO product_categories (product_id, category_id) VALUES (?, ?)",
			productID, cid,
		).Error; err != nil {
			return err
		}
	}
	return nil
}