package category

import "website-api/model/category"

func (r *repo) Create(cat category.Category) error {
	if cat.ParentId == "" {
		return r.db.Exec(
			"INSERT INTO categories (id, name, slug, parent_id, created_at, updated_at) VALUES (?, ?, ?, NULL, NOW(), NULL)",
			cat.Id, cat.Name, cat.Slug,
		).Error
	}
	return r.db.Omit("created_at", "updated_at").Create(&cat).Error
}

func (r *repo) FindByID(id string) (cat category.Category, err error) {
	return cat, r.db.Where("id = ?", id).First(&cat).Error
}

func (r *repo) FindBySlug(slug string) (cat category.Category, err error) {
	return cat, r.db.Where("slug = ?", slug).First(&cat).Error
}

func (r *repo) Update(id string, values map[string]any) error {
	return r.db.Model(&category.Category{}).Where("id = ?", id).Updates(values).Error
}

func (r *repo) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&category.Category{}).Error
}