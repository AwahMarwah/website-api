package auth

import (
	authModel "website-api/model/auth"
	authRepo "website-api/repository/auth"
	userRepo "website-api/repository/user"
)

type (
	IService interface {
		ForgotPassword(req *authModel.ForgotPasswordRequest) (statusCode int, message string, err error)
		ResetPassword(req *authModel.ResetPasswordRequest) (statusCode int, message string, err error)
		ResendVerification(req *authModel.ResendVerificationRequest) (statusCode int, err error)
		Authorize(authHeader *string) (userID string, statusCode int, err error)
	}

	service struct {
		authRepo authRepo.IRepo
		userRepo userRepo.IRepo
	}
)

func NewService(userRepo userRepo.IRepo, authRepo authRepo.IRepo) IService {
	return &service{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}
