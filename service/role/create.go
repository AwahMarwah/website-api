package role

import (
	"errors"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
	roleModel "website-api/model/role"
)

func (s *service) Create(reqBody *roleModel.RoleCreateReqBody) (statusCode int, err error) {
	role, err := s.roleRepo.Take([]string{"id", "name"}, &roleModel.Role{Name: reqBody.Name + " AND deleted_at IS NULL"})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusInternalServerError, fmt.Errorf("gagal memeriksa role yang sudah ada: %w", err)
	}
	if role.Name != "" {
		return http.StatusConflict, fmt.Errorf("role dengan nama '%s' sudah ada", reqBody.Name)
	}
	return http.StatusCreated, s.roleRepo.Create(&roleModel.Role{
		Name:        strings.ToLower(reqBody.Name),
		DisplayName: cases.Title(language.Indonesian).String(strings.ToLower(reqBody.DisplayName)),
		Description: reqBody.Description,
		IsActive:    reqBody.IsActive,
		CreatedAt:   time.Now(),
		CreatedBy:   "system",
	})
}
