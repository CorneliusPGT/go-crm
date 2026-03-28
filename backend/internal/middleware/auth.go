package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				http.Error(w, "Не авторизирован", http.StatusUnauthorized)
				return
			}
			claims := &Claims{}
			tokenStr := strings.TrimPrefix(header, "Bearer ")
			token, err := jwt.ParseWithClaims(
				tokenStr,
				claims,
				func(t *jwt.Token) (interface{}, error) {
					if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Неверный алгоритм")
					}
					return []byte(secret), nil
				},
			)
			if err != nil || !token.Valid {
				http.Error(w, "Не авторизирован", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
