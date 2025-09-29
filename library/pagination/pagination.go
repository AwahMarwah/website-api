package pagination

func Offset(limit, page *int) int {
	if *limit == 0 {
		*limit = 10
	}
	if *page == 0 {
		*page = 1
	}
	return (*page - 1) * *limit
}
