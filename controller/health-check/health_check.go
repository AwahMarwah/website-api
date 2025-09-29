package health_check

import (
	"database/sql"
	health_check_repo "website-api/repository/health-check"
	health_check "website-api/service/health-check"
)

type controller struct {
	healthService health_check.IService
}

func NewController(db *sql.DB) *controller {
	return &controller{
		healthService: health_check.NewService(health_check_repo.NewRepo(db)),
	}
}
