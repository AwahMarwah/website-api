package content_page

import (
	"fmt"
	"time"
	libStruct "website-api/library"
	lib "website-api/library/cache"
	modelContentPage "website-api/model/content-page"
)

func (s *service) FaqList(reqQuery *modelContentPage.FaqListReqQuery) (resData []modelContentPage.FaqListResponse, count int64, err error) {
	cacheKey := lib.GenerateCacheKey(libStruct.GetStructName(modelContentPage.CmsFaq{}), reqQuery)

	err = s.redis.Get(cacheKey, &resData)
	if err == nil {
		fmt.Printf("CACHE HIT | key=%s\n", cacheKey)
		return resData, count, nil
	}

	fmt.Printf("CACHE MISS | key=%s\n", cacheKey)

	resData, count, err = s.contentPageRepo.FaqList(reqQuery)
	if err != nil {
		return nil, count, err
	}

	_ = s.redis.Set(cacheKey, resData, 1*time.Hour)
	fmt.Printf("CACHE SET | key=%s | ttl=5m\n", cacheKey)

	return resData, count, nil
}
