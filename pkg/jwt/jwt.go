package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var signingMethod = jwt.SigningMethodHS256

func BuildExpirable(dur time.Duration, key string) (string, error) {
	token := jwt.New(signingMethod)

	token.Claims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(dur)),
	}

	return token.SignedString([]byte(key))
}

func Validate(tokenString string, key string) error {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected sign method: %v", t.Header["alg"])
		}

		return []byte(key), nil
	})
	if err != nil {
		return fmt.Errorf("jwt: %w", err)
	}

	if !token.Valid {
		return fmt.Errorf("token not valid")
	}

	// if time.Since(token.Claims.ExpiresAt.Time) > 0 {
	// 	return c.HelloMessage, fmt.Errorf("token expired")
	// }

	return nil
}
