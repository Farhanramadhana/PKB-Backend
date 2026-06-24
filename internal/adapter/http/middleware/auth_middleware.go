package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/bpka/tps-pkb/internal/adapter/http/response"
	"github.com/bpka/tps-pkb/internal/domain/port"
)

type contextKey string

const (
	CtxKeyClaims    contextKey = "claims"
	CtxKeyIPAddress contextKey = "ip_address"
	CtxKeyUserAgent contextKey = "user_agent"
)

func ClaimsFromContext(ctx context.Context) *port.Claims {
	c, _ := ctx.Value(CtxKeyClaims).(*port.Claims)
	return c
}

func IPFromContext(ctx context.Context) string {
	s, _ := ctx.Value(CtxKeyIPAddress).(string)
	return s
}

func UserAgentFromContext(ctx context.Context) string {
	s, _ := ctx.Value(CtxKeyUserAgent).(string)
	return s
}

func Auth(tokenProvider port.TokenProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				response.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "token diperlukan")
				return
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := tokenProvider.Validate(tokenStr)
			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "token tidak valid")
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, CtxKeyClaims, claims)
			ctx = context.WithValue(ctx, CtxKeyIPAddress, realIP(r))
			ctx = context.WithValue(ctx, CtxKeyUserAgent, r.Header.Get("User-Agent"))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func realIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
