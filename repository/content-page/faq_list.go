package content_page

import (
	"fmt"

	modelContentPage "website-api/model/content-page"
)

func (r *repo) FaqList(reqQuerry *modelContentPage.FaqListReqQuery) (resData []modelContentPage.FaqListResponse, count int64, err error) {
	resData = make([]modelContentPage.FaqListResponse, 0)
	query := r.db.Model(&modelContentPage.CmsFaq{}).Where("is_active = ?", true)

	if reqQuerry.Search != "" {
		pattern := fmt.Sprintf("%%%s%%", reqQuerry.Search)
		query = query.Where("question ILIKE ? OR answer ILIKE ?", pattern, pattern)
	}

	if err = query.Count(&count).Error; err != nil {
		return
	}
	err = query.Order("order_no ASC").Limit(reqQuerry.Limit).Offset(reqQuerry.Offset).Scan(&resData).Error
	return
}

func (r *repo) FaqListAdmin(reqQuerry *modelContentPage.FaqListReqQuery) (resData []modelContentPage.CmsFaq, count int64, err error) {
	resData = make([]modelContentPage.CmsFaq, 0)
	query := r.db.Model(&modelContentPage.CmsFaq{})

	if reqQuerry.Search != "" {
		pattern := fmt.Sprintf("%%%s%%", reqQuerry.Search)
		query = query.Where("question ILIKE ? OR answer ILIKE ?", pattern, pattern)
	}

	if err = query.Count(&count).Error; err != nil {
		return
	}
	err = query.Order("order_no ASC").Limit(reqQuerry.Limit).Offset(reqQuerry.Offset).Find(&resData).Error
	return
}

func (r *repo) FaqFindByID(id string) (modelContentPage.CmsFaq, error) {
	var faq modelContentPage.CmsFaq
	err := r.db.Where("id = ?", id).First(&faq).Error
	return faq, err
}

func (r *repo) FaqCreate(faq *modelContentPage.CmsFaq) error {
	return r.db.Create(faq).Error
}

func (r *repo) FaqUpdate(id string, values map[string]any) error {
	return r.db.Model(&modelContentPage.CmsFaq{}).Where("id = ?", id).Updates(values).Error
}

func (r *repo) FaqDelete(id string) error {
	return r.db.Where("id = ?", id).Delete(&modelContentPage.CmsFaq{}).Error
}