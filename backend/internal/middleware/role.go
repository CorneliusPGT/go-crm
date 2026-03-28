package middleware

import "net/http"

func AdminMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey).(*Claims)
		if !ok || claims.Role != "admin" {
			http.Error(w, "Запрещено", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}
