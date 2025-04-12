package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/spf13/viper"
)

var (
	ErrInvalidTokenExpired = errx.InternalServerError("invalid token expired")
	ErrFailedClaimJWT      = errx.InternalServerError("failed claim jwt")
	ErrInvalidSignature    = errx.InternalServerError("invalid signature")
	ErrSignJwt             = errx.InternalServerError("failed to sign jwt")
)

func EncodeToken(account *entity.Account) (string, error) {
	claims := &JWTClaims{
		Role: account.Role.Name,
		ID:   account.Id.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(viper.GetDuration("jwt.expiration"))),
			Issuer:    viper.GetString("jwt.issuer"),
			Subject:   account.Id.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(viper.GetString("JWT_SECRET_KEY")))
	if err != nil {
		return "", ErrSignJwt
	}
	return signedToken, nil
}

func DecodeToken(token string) (*JWTClaims, error) {
	decoded, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		token.Method = jwt.SigningMethodHS256
		return []byte(viper.GetString("jwt.secretkey")), nil
	})

	if err != nil {
		return &JWTClaims{}, ErrInvalidSignature
	}

	if !decoded.Valid {
		return &JWTClaims{}, ErrInvalidTokenExpired
	}

	claims, ok := decoded.Claims.(*JWTClaims)
	if !ok {
		return &JWTClaims{}, ErrFailedClaimJWT
	}

	return claims, nil
}
