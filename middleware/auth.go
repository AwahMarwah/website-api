package middleware

import (
	"errors"
	"net/http"
	"strings"
	"website-api/database/encrypt"
	"website-api/library/response"
	authModel "website-api/model/auth"
	authRepo "website-api/repository/auth"
	userRepo "website-api/repository/user"
	"website-api/service/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserInfo struct {
	UserID   string
	Email    string
	Username string
	Role     string
}

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqHeader authModel.ReqHeader
		if err := ctx.ShouldBindHeader(&reqHeader); err != nil {
			response.Error(ctx, http.StatusUnauthorized, err.Error())
			ctx.Abort()
			return
		}

		tokenParts := strings.Split(reqHeader.Authorization, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			response.Error(ctx, http.StatusUnauthorized, "Format token tidak valid")
			ctx.Abort()
			return
		}

		tokenString := tokenParts[1]

		_, claims, err := encrypt.Parse("Bearer " + tokenString)
		if err != nil {
			response.Error(ctx, http.StatusUnauthorized, "Token tidak valid")
			ctx.Abort()
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			response.Error(ctx, http.StatusUnauthorized, "Token tidak valid")
			ctx.Abort()
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			response.Error(ctx, http.StatusUnauthorized, "Token tidak valid")
			ctx.Abort()
			return
		}
		roleName, ok := claims["role_name"].(string)
		if !ok {
			response.Error(ctx, http.StatusUnauthorized, "Token tidak valid")
			ctx.Abort()
			return
		}

		authService := auth.NewService(userRepo.NewRepo(db), authRepo.NewRepo(db))
		userId, statusCode, err := authService.Authorize(&reqHeader.Authorization)
		if err != nil {
			response.Error(ctx, statusCode, err.Error())
			ctx.Abort()
			return
		}

		ctx.Set("role_name", roleName)
		ctx.Set("user_id", userId)
		ctx.Set("user_email", email)
		ctx.Set("username", username)
		ctx.Next()
	}
}

// GetUserFromContext mengambil informasi user yang sudah terautentikasi dari context
func GetUserFromContext(ctx *gin.Context) (*UserInfo, error) {
	userID, ok := ctx.Get("user_id")
	if !ok {
		return nil, errors.New("user ID tidak ditemukan dalam context")
	}

	email, ok := ctx.Get("user_email")
	if !ok {
		return nil, errors.New("email tidak ditemukan dalam context")
	}

	username, ok := ctx.Get("username")
	if !ok {
		return nil, errors.New("username tidak ditemukan dalam context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return nil, errors.New("format user ID tidak valid")
	}

	emailStr, ok := email.(string)
	if !ok {
		return nil, errors.New("format email tidak valid")
	}

	usernameStr, ok := username.(string)
	if !ok {
		return nil, errors.New("format username tidak valid")
	}

	role, ok := ctx.Get("role_name")
	if !ok {
		return nil, errors.New("role tidak ditemukan dalam context")
	}

	roleStr, ok := role.(string)
	if !ok {
		return nil, errors.New("format role tidak valid")
	}
	return &UserInfo{
		UserID:   userIDStr,
		Email:    emailStr,
		Username: usernameStr,
		Role:     roleStr,
	}, nil
}
