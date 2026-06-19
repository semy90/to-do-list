package jwt_utils

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Id int
	*jwt.RegisteredClaims
}

func GetIdFromToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("SALT")), nil
	})
	if err != nil {
		return -1, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.Id, nil
	} else {
		return -1, fmt.Errorf("error with claims in jwt")
	}
}
