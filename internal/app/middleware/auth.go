package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/nontypeable/financial-tracker/internal/auth"
	contextKeys "github.com/nontypeable/financial-tracker/internal/context"
	httpHelper "github.com/nontypeable/financial-tracker/internal/http"
)

func AuthMiddleware(tokenManager auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				httpHelper.Error(w, http.StatusUnauthorized, "authorization header missing or malformed")
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := tokenManager.ValidateAccessToken(token)
			if err != nil {
				w.Header().Set("X-Token-Expired", "true")
				httpHelper.Error(w, http.StatusUnauthorized, "invalid or expired access token")
				return
			}

			ctx := context.WithValue(r.Context(), contextKeys.UserIDKey, claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
