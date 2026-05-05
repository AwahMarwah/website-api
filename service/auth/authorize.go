package auth

import (
	"errors"
	"net/http"
	"website-api/common"
	"website-api/database/encrypt"
	modelUser "website-api/model/user"

	"gorm.io/gorm"
)

func (s *service) Authorize(token *string) (userId string, statusCode int, err error) {
	tokenRaw, claims, err := encrypt.Parse(*token)
	if err != nil {
		return userId, http.StatusUnauthorized, err
	}
	userId = claims["user_id"].(string)
	user, err := s.userRepo.Take([]string{"token"}, &modelUser.User{Id: userId})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userId, http.StatusUnauthorized, errors.New(common.UserNotFound)
		}
		return userId, http.StatusInternalServerError, err
	}
	if tokenRaw != user.Token {
		return userId, http.StatusUnauthorized, errors.New(common.UserHasSignedOut)
	}
	return userId, http.StatusOK, nil
}
