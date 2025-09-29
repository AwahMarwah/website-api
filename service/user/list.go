package user

import (
	userModel "website-api/model/user"
)

func (s *service) List(reqQuery *userModel.ListUserReqQuery) (resData []userModel.ListUserResponse, count int64, err error) {
	resData = make([]userModel.ListUserResponse, 0)
	user, count, err := s.userRepo.Find(reqQuery)
	if err != nil {
		return nil, 0, err
	}

	for i := range user {
		var (
			createdAtFormatted string
			updatedAtFormatted string
			deletedAtFormatted string
		)
		if !user[i].CreatedAt.IsZero() {
			createdAtFormatted = user[i].CreatedAt.Format("2006-Jan-02")
		}

		// Format UpdatedAt jika ada
		if user[i].UpdatedAt != nil && !user[i].UpdatedAt.IsZero() {
			updatedAtFormatted = user[i].UpdatedAt.Format("2006-Jan-02")
		}

		// Format DeletedAt jika ada
		if user[i].DeletedAt != nil && !user[i].DeletedAt.IsZero() {
			deletedAtFormatted = user[i].DeletedAt.Format("2006-Jan-02")
		}
		resData = append(resData, userModel.ListUserResponse{
			Id:                 user[i].Id,
			Name:               user[i].Name,
			UserName:           user[i].UserName,
			Picture:            user[i].Picture,
			Email:              user[i].Email,
			PhoneNumber:        user[i].PhoneNumber,
			IsVerified:         user[i].IsVerified,
			CreatedAt:          user[i].CreatedAt,
			CreatedAtFormatted: createdAtFormatted,
			CreatedBy:          user[i].CreatedBy,
			UpdatedAt:          user[i].UpdatedAt,
			UpdatedByFormatted: updatedAtFormatted,
			UpdatedBy:          user[i].UpdatedBy,
			DeletedAt:          user[i].DeletedAt,
			DeletedByFormatted: deletedAtFormatted,
			DeletedBy:          user[i].DeletedBy,
		})
	}
	return
}
