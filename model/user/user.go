package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id                         string
	Name                       string
	UserName                   string
	Picture                    string
	Email                      string
	EncryptedPassword          string
	Token                      string
	PhoneNumber                string
	VerificationToken          string
	VerificationTokenExpiredAt time.Time
	IsVerified                 bool
	CreatedAt                  time.Time
	CreatedBy                  string
	UpdatedAt                  *time.Time
	UpdatedBy                  string
	DeletedAt                  *time.Time
	DeletedBy                  string
	RoleId                     string
}

func (user *User) BeforeCreate(*gorm.DB) error {
	user.Id = uuid.New().String()
	return nil
}
