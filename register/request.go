package register

import "github.com/dgrijalva/jwt-go"

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	UserID uint `json:"id"`
	jwt.StandardClaims
}
