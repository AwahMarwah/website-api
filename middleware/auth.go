package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"website-api/database/encrypt"
	"website-api/library/response"
)

type UserInfo struct {
	UserID   string
	Email    string
	Username string
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(ctx, http.StatusUnauthorized, "Token tidak ditemukan")
			ctx.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
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

		userID, ok := claims["user_id"].(string)
		if !ok {
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

		ctx.Set("user_id", userID)
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

	return &UserInfo{
		UserID:   userIDStr,
		Email:    emailStr,
		Username: usernameStr,
	}, nil
}
