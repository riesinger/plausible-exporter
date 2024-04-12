package server

import (
	"net/http"
)

func BearerAuthMiddleware(next http.Handler, token string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		headerValue := r.Header.Get("Authorization")
		if headerValue == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		expectedToken := "Bearer " + token
		if headerValue != expectedToken {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
