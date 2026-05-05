package user

func (s *service) SignOut(userId string) (err error) {
	return s.userRepo.Update(&userId, &map[string]any{"token": ""})
}
