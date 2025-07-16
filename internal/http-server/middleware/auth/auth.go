package auth

import (
	"net/http"
	"strings"
	"vkTestMarketplace/internal/http-server/handlers"
	"vkTestMarketplace/internal/lib/jwt"
	utilsauth "vkTestMarketplace/internal/utils/auth"
)

type AdHandler struct {
	Handlers *handlers.AdHandlers
}

func (h *AdHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "missing auth header", http.StatusUnauthorized)
			return
		} else if !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, "missing or invalid token", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		userID, err := jwt.ParseToken(token, h.Handlers.JWTSecret)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = utilsauth.WithUserID(ctx, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OptionalAuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header != "" && strings.HasPrefix(header, "Bearer ") {
				token := strings.TrimPrefix(header, "Bearer ")
				userID, err := jwt.ParseToken(token, secret)
				if err == nil {
					ctx := utilsauth.WithUserID(r.Context(), userID)
					r = r.WithContext(ctx)
				} else {
					http.Error(w, "invalid token", http.StatusUnauthorized)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
