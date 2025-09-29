package role

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"time"
	modelRole "website-api/model/role"
)

func (s *service) Update(req *modelRole.UpdateReq) (statusCode int, err error) {
	if _, err = s.roleRepo.Take([]string{"id"}, &modelRole.Role{Id: req.Path.Id}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}
	values := map[string]any{
		"name":         req.Body.Name,
		"display_name": req.Body.DisplayName,
		"description":  req.Body.Description,
		"is_active":    req.Body.IsActive,
		"updated_at":   time.Now(),
	}
	if err = s.roleRepo.Update(&req.Path.Id, &values); err != nil {
		return http.StatusInternalServerError, err
	}
	return
}
