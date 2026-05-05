package user

import (
	"errors"
	"fmt"
	"net/http"
	roleModel "website-api/model/role"
	userModel "website-api/model/user"

	"gorm.io/gorm"
)

func (s *service) Update(req *userModel.UserUpdateRequest) (statusCode int, err error) {
	user, err := s.userRepo.Take([]string{"id"}, &userModel.User{Id: req.Path.Id})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusInternalServerError, fmt.Errorf("failed to check existing user: %w", err)
	}
	if user.Id == "" {
		return http.StatusBadRequest, fmt.Errorf("user %s tidak ditemukan", req.Path.Id)
	}

	role, err := s.roleRepo.Take([]string{"id"}, &roleModel.Role{Id: req.Body.RoleID})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusInternalServerError, fmt.Errorf("failed to check role: %w", err)
	}

	if role.Id == "" {
		return http.StatusBadRequest, fmt.Errorf("role %s tidak ditemukan", req.Path.Id)
	}

	values := map[string]any{
		"role_id":      req.Body.RoleID,
		"phone_number": req.Body.PhoneNumber,
	}
	if err = s.userRepo.Update(&req.Path.Id, &values); err != nil {
	}
	return
}
