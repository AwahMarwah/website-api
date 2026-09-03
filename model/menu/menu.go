package menu

import "time"

type Menu struct {
	ID          string
	ParentID    *string
	Name        string
	DisplayName string
	Icon        string
	Path        string
	SortOrder   int
	IsActive    bool
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   *time.Time
	UpdatedBy   string
	DeletedAt   *time.Time
	DeletedBy   string
}
