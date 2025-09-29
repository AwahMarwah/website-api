package role

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"time"
	modelRole "website-api/model/role"
)

func (s *service) Delete(reqPath *modelRole.ReqPath) (statusCode int, err error) {
	if _, err = s.roleRepo.Take([]string{"id"}, &modelRole.Role{Id: reqPath.Id}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}
	values := map[string]any{
		"deleted_at": time.Now(),
		"deleted_by": "system",
	}
	if err = s.roleRepo.Update(&reqPath.Id, &values); err != nil {
		return http.StatusInternalServerError, err
	}
	return
}
