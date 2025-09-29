package health_check

func (s *service) Check() (err error) {
	return s.healthRepo.Ping()
}
