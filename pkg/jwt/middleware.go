package auth

import (
	"errors"
	"net/http"
	"strings"
)

func WithJWT(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, err := GetBearer(r)
			if err != nil {
				http.Error(w, "jwt expected", http.StatusUnauthorized)
				return
			}

			err = Validate(b, key)
			if err != nil {
				http.Error(w, "jwt invalid", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetBearer(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	b, found := strings.CutPrefix(h, "Bearer ")
	if !found {
		return "", errors.New("bearer not found")
	}

	return b, nil
}
