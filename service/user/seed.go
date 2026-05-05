package user

func (s *service) Seed() (err error) {
	return s.userRepo.Seed()
}
