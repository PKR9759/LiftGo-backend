// internal/auth/middleware.go
package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/PKR9759/LiftGo-backend/internal/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

// RequireAuth protects any route that needs a logged-in user
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			utils.WriteError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.WriteError(w, http.StatusUnauthorized, "invalid authorization format")
			return
		}

		claims, err := ValidateToken(parts[1])
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// attach claims to request context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext pulls the logged-in user's claims out of context
func GetUserFromContext(r *http.Request) *Claims {
	claims, _ := r.Context().Value(UserContextKey).(*Claims) //type casting into *Claims	
	return claims
}