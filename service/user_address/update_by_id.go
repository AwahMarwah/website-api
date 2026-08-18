package user_address

import (
	"errors"
	"fmt"
	userAddressModel "website-api/model/user_address"
)

func (s *service) UpdateByID(req *userAddressModel.ReqUpdateUserAddress) error {
	// check existing data
	userAddress, err := s.userAddressRepo.Take([]string{"id"}, &userAddressModel.UserAddress{ID: req.Path.ID})
	if err != nil {
		return err
	}
	fmt.Println(userAddress, "service address")

	if userAddress.ID == "" {
		return errors.New("record not found")
	}

	values := map[string]any{
		"recipient_name": req.Body.RecepientName,
		"city":           req.Body.City,
		"postal_code":    req.Body.PostalCode,
		"phone_number":   req.Body.PhoneNumber,
		"full_address":   req.Body.FullAddress,
	}

	return s.userAddressRepo.Update(&req.Path.ID, &values)
}
