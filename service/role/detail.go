package role

import (
	"net/http"
	modelRole "website-api/model/role"
)

func (s *service) Detail(reqPath *modelRole.ReqPath) (resData modelRole.DetailResponse, statusCode int, err error) {
	roleData, err := s.roleRepo.Take([]string{"id", "name", "display_name", "description", "is_active", "created_at", "updated_at"}, &modelRole.Role{Id: reqPath.Id})
	if err != nil {
		return resData, http.StatusInternalServerError, err
	}
	resData = modelRole.DetailResponse{
		Id:          roleData.Id,
		Name:        roleData.Name,
		DisplayName: roleData.DisplayName,
		Description: roleData.Description,
		IsActive:    roleData.IsActive,
		CreatedAt:   roleData.CreatedAt,
		UpdatedAt:   roleData.UpdatedAt,
	}
	return resData, statusCode, nil
}
