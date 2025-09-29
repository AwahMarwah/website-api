package health_check

import "database/sql"

type (
	IRepo interface {
		Ping() (err error)
	}

	repo struct {
		db *sql.DB
	}
)

func NewRepo(db *sql.DB) IRepo {
	return &repo{db: db}
}
