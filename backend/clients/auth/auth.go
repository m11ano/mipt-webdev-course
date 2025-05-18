package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type ClientImpl struct {
	jwtSecretKey string
}

func NewClient(jwtSecretKey string) *ClientImpl {
	return &ClientImpl{jwtSecretKey: jwtSecretKey}
}

func (a *ClientImpl) ParseJWT(tokenStr string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AuthClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(a.jwtSecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
