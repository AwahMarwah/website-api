package user_address

import (
	"time"
	libStruct "website-api/library"
	lib "website-api/library/cache"
	"website-api/model/user_address"
)

func (s *service) Find(userID string) (resData []user_address.UserAddress, err error) {
	cacheKey := lib.GenerateCacheKey(libStruct.GetStructName(user_address.UserAddress{}), userID)

	err = s.redis.Get(cacheKey, &resData)
	if err == nil {
		return resData, nil
	}
	userAddress, err := s.userAddressRepo.Find(userID)
	if err != nil {
		return nil, err
	}

	resData = userAddress

	_ = s.redis.Set(cacheKey, resData, 1*time.Hour)

	return
}
