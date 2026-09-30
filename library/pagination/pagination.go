package pagination

func Offset(limit, page *int) int {
	if *limit <= 0 {
		*limit = 10
	}
	if *limit > 100 {
		*limit = 100
	}
	if *page <= 0 {
		*page = 1
	}
	return (*page - 1) * *limit
}