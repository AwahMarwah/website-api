package brand

import brandModel "website-api/model/brand"

func (r *repo) Create(brandData brandModel.Brand) error {
	return r.db.Omit("created_at", "updated_at").Create(&brandData).Error
}

func (r *repo) FindByID(id string) (brand brandModel.Brand, err error) {
	return brand, r.db.Where("id = ?", id).First(&brand).Error
}

func (r *repo) Update(id string, values map[string]any) error {
	return r.db.Model(&brandModel.Brand{}).Where("id = ?", id).Updates(values).Error
}

func (r *repo) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&brandModel.Brand{}).Error
}