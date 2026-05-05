package content_page

import (
	modelContentPage "website-api/model/content-page"
)

func (r *repo) Take(selectParams []string, conditions *modelContentPage.CmsPage) (cmsPage modelContentPage.CmsPage, err error) {
	return cmsPage, r.db.Debug().Select(selectParams).Take(&cmsPage, conditions).Error
}
