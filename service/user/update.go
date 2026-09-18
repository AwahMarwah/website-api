package user

import (
	"errors"
	"fmt"
	"net/http"
	"time"
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

	values := map[string]any{}

	// Update phone jika dikirim
	if req.Body.PhoneNumber != "" {
		values["phone_number"] = req.Body.PhoneNumber
	}

	// Update role jika dikirim
	if req.Body.RoleID != "" {
		role, roleErr := s.roleRepo.Take([]string{"id"}, &roleModel.Role{Id: req.Body.RoleID})
		if roleErr != nil && !errors.Is(roleErr, gorm.ErrRecordNotFound) {
			return http.StatusInternalServerError, fmt.Errorf("failed to check role: %w", roleErr)
		}
		if role.Id == "" {
			return http.StatusBadRequest, fmt.Errorf("role %s tidak ditemukan", req.Body.RoleID)
		}
		values["role_id"] = req.Body.RoleID
	}

	if len(values) == 0 {
		return http.StatusBadRequest, fmt.Errorf("tidak ada data yang dikirim")
	}
	values["updated_at"] = time.Now()

	if err = s.userRepo.Update(&req.Path.Id, &values); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to update user: %w", err)
	}
	return http.StatusOK, nil
}
