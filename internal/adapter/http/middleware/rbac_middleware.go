package middleware

import (
	"context"
	"net/http"

	"github.com/bpka/tps-pkb/internal/adapter/http/response"
	"github.com/bpka/tps-pkb/internal/domain/entity"
	"github.com/google/uuid"
)

type Middleware func(http.Handler) http.Handler

type contextOwnerKey string

const CtxKeyOwnerFilter contextOwnerKey = "owner_filter"

func OwnerFilterFromContext(ctx context.Context) *uuid.UUID {
	v, _ := ctx.Value(CtxKeyOwnerFilter).(*uuid.UUID)
	return v
}

// Chain applies middleware right-to-left so the first listed runs outermost.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// RequireRole rejects requests whose JWT role is not in the allowed set.
func RequireRole(allowed ...entity.Role) Middleware {
	set := make(map[entity.Role]struct{}, len(allowed))
	for _, r := range allowed {
		set[r] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := ClaimsFromContext(r.Context())
			if claims == nil {
				response.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "tidak terautentikasi")
				return
			}
			if _, ok := set[claims.Role]; !ok {
				response.WriteError(w, http.StatusForbidden, "FORBIDDEN", "akses ditolak")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// SelfOnly injects an owner filter for WAJIB_PAJAK role so they can only see their own data.
// For other roles it is a no-op.
func SelfOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := ClaimsFromContext(r.Context())
		if claims != nil && claims.Role == entity.RoleWajibPajak {
			if claims.WajibPajakID == nil {
				response.WriteError(w, http.StatusForbidden, "FORBIDDEN", "akun tidak terhubung ke wajib pajak")
				return
			}
			ctx := context.WithValue(r.Context(), CtxKeyOwnerFilter, claims.WajibPajakID)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
