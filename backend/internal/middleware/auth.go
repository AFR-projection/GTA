package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/AFR-projection/GTA/backend/internal/auth"
	"github.com/AFR-projection/GTA/backend/internal/httpx"
	"github.com/google/uuid"
)

type ctxKey string

const accountIDKey ctxKey = "account_id"

func AccountIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(accountIDKey).(uuid.UUID)
	return id, ok
}

func Auth(tokens *auth.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				httpx.Error(w, http.StatusUnauthorized, "missing bearer token")
				return
			}
			raw := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
			claims, err := tokens.Parse(raw)
			if err != nil {
				httpx.Error(w, http.StatusUnauthorized, "invalid token")
				return
			}
			ctx := context.WithValue(r.Context(), accountIDKey, claims.AccountID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
