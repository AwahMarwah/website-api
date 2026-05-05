package content_page

func (s *service) Seed() (err error) {
	return s.contentPageRepo.SeedCmsPage()
}
