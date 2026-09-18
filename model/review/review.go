package review

import "time"

type Review struct {
	ID        string
	ProductID string
	UserID    string
	Rating    int
	Comment   string
	CreatedAt time.Time
}