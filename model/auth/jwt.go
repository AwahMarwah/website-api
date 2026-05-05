package auth

import "github.com/dgrijalva/jwt-go"

type JWTClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	RoleName string `json:"role_name"`
	jwt.StandardClaims
}
