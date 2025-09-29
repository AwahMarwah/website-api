package role

func (s *service) Seed() (err error) {
	return s.roleRepo.Seed()
}
