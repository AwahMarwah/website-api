package user

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	roleModel "website-api/model/role"
	userModel "website-api/model/user"
)

func (s *service) Detail(reqPath *userModel.ReqPath) (resData userModel.DetailResponse, statusCode int, err error) {
	user, err := s.userRepo.Take([]string{"id", "name", "user_name", "email", "encrypted_password", "phone_number", "is_verified", "role_id"}, &userModel.User{Id: reqPath.Id})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resData, http.StatusNotFound, fmt.Errorf("user not found")
		}
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal mengambil data user: %w", err)
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
	}
	return resData, statusCode, nil
}
