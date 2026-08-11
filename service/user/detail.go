package user

import (
	"errors"
	"fmt"
	"net/http"
	roleModel "website-api/model/role"
	userModel "website-api/model/user"
	"website-api/model/user_address"

	"gorm.io/gorm"
)

func (s *service) Detail(reqPath *userModel.ReqPath) (resData userModel.DetailResponse, statusCode int, err error) {
	user, err := s.userRepo.Take([]string{"id", "name", "user_name", "email", "encrypted_password", "phone_number", "is_verified", "role_id"}, &userModel.User{Id: reqPath.Id})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resData, http.StatusNotFound, fmt.Errorf("user not found")
		}
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal mengambil data user: %w", err)
	}
	userAddress, err := s.userAddressRepo.Take([]string{"id", "user_id", "recipient_name", "phone_number", "full_address", "city", "postal_code", "is_primary", "created_at", "updated_at"}, &user_address.UserAddress{UserID: reqPath.Id, IsPrimary: true})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resData, http.StatusNotFound, fmt.Errorf("user address not found")
		}
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal mengambil data user address: %w", err)
	}
	role, err := s.roleRepo.Take([]string{"id", "name", "display_name", "description"}, &roleModel.Role{Id: user.RoleId})
	if err != nil {
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal mengambil data role: %w", err)
	}
	resData = userModel.DetailResponse{
		ID:          user.Id,
		Name:        user.Name,
		UserName:    user.UserName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		IsVerified:  user.IsVerified,
		Role:        role,
		Address:     userAddress,
	}
	return resData, statusCode, nil
}
