package health_check

func (r *repo) Ping() (err error) {
	return r.db.Ping()
}
