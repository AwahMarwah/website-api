package user_address

import (
	userAddressModel "website-api/model/user_address"
)

func (s *service) Create(reqBody *userAddressModel.CreateUserAddressRequest) error {
	userAddress := userAddressModel.UserAddress{
		UserID:        reqBody.UserID,
		IsPrimary:     reqBody.IsPrimary,
		PhoneNumber:   reqBody.PhoneNumber,
		City:          reqBody.City,
		FullAddress:   reqBody.FullAddress,
		PostalCode:    reqBody.PostalCode,
		RecipientName: reqBody.RecipientName,
	}
	return s.userAddressRepo.Create(&userAddress)
}
