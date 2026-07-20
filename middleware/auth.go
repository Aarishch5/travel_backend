package middleware

import (
	"context"
	"net/http"
	"strings"

	"TravelBackend/utils"
)

type contextKey string

const ClaimsContextKey contextKey = "claims"

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.RespondError(w, http.StatusUnauthorized, "authorization header must be a bearer token")
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := ClaimsFromContext(r.Context())
		if !ok {
			utils.RespondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		if claims.Role != role {
			utils.RespondError(w, http.StatusForbidden, "forbidden")
			return
		}
		next.ServeHTTP(w, r)
	}
}

func ClaimsFromContext(ctx context.Context) (*utils.Claims, bool) {
	claims, ok := ctx.Value(ClaimsContextKey).(*utils.Claims)
	if !ok {
		return nil, false
	}
	return claims, ok
}
