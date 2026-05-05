package user

import (
	modelUser "website-api/model/user"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		Create(user *modelUser.User) (err error)
		Find(reqQuery *modelUser.ListUserReqQuery) (users []modelUser.ListUserResponse, count int64, err error)
		Take(selectParams []string, conditions *modelUser.User) (user modelUser.User, err error)
		Update(id *string, values *map[string]any) (err error)
		Seed() (err error)
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}
