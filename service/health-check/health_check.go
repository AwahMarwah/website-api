package health_check

import health_check "website-api/repository/health-check"

type (
	IService interface {
		Check() (err error)
	}

	service struct {
		healthRepo health_check.IRepo
	}
)

func NewService(healthRepo health_check.IRepo) IService {
	return &service{
		healthRepo: healthRepo,
	}
}
