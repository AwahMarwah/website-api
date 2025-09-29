package auth

import (
	"gorm.io/gorm"
	authModel "website-api/model/auth"
)

type (
	IRepo interface {
		Create(auth *authModel.PasswordResetToken) (err error)
		Take(selectParams []string, condition *authModel.PasswordResetToken) (result authModel.PasswordResetToken, err error)
		Update(id *string, values *map[string]any) (err error)
		Delete(condition map[string]interface{}) (err error)
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}
