package jwt

import "github.com/golang-jwt/jwt/v4"

type JWTClaims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}
