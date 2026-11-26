package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/alexander-ds/EjercicioIA-inventario/inventory-service-go/internal/auth"
)

type contextKey string

const claimsKey contextKey = "claims"

// requireAuth rechaza peticiones sin Bearer JWT válido con 401.
func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Fields(r.Header.Get("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			writeError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		claims, err := s.validator.Parse(parts[1])
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// claimsFrom permite a futuros handlers obtener el user_id del claim `sub`.
func claimsFrom(ctx context.Context) *auth.Claims {
	claims, _ := ctx.Value(claimsKey).(*auth.Claims)
	return claims
}
