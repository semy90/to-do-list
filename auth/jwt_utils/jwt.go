package jwt_utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Id int
	*jwt.RegisteredClaims
}

// return access token, error
func NewAccessToken(id int) (string, error) {
	SignedKey := []byte(os.Getenv("SALT"))
	claims := CustomClaims{
		id,
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(SignedKey)
	return ss, err
}

// return access token, refresh token, error
func NewPairOfTokens(id int) (string, string, error) {
	SignedKey := []byte(os.Getenv("SALT"))
	claimsAccess := CustomClaims{
		id,
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	claimsRefresh := CustomClaims{
		id,
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	refreshString, err := refreshToken.SignedString(SignedKey)
	if err != nil {
		return "", "", err
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccess)
	accessString, err := accessToken.SignedString(SignedKey)
	if err != nil {
		return "", "", err
	}
	return accessString, refreshString, nil
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
