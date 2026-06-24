package http

import (
	"net/http"

	"github.com/bpka/tps-pkb/internal/adapter/http/handler"
	"github.com/bpka/tps-pkb/internal/adapter/http/middleware"
	"github.com/bpka/tps-pkb/internal/domain/entity"
	"github.com/bpka/tps-pkb/internal/domain/port"
)

type Handlers struct {
	Auth       *handler.AuthHandler
	WajibPajak *handler.WajibPajakHandler
	Kewajiban  *handler.KewajibanHandler
	Pembayaran *handler.PembayaranHandler
	Denda      *handler.DendaHandler
	Laporan    *handler.LaporanHandler
	Health     *handler.HealthHandler
	Docs       *handler.DocsHandler
}

func NewRouter(h *Handlers, tokenProvider port.TokenProvider) http.Handler {
	mux := http.NewServeMux()

	authMW := middleware.Auth(tokenProvider)
	admin := middleware.RequireRole(entity.RoleAdmin)
	staff := middleware.RequireRole(entity.RoleAdmin, entity.RolePetugas)
	all := middleware.RequireRole(entity.RoleAdmin, entity.RolePetugas, entity.RoleWajibPajak)

	chain := middleware.Chain

	// Health check (no auth)
	mux.HandleFunc("GET /health", h.Health.Check)

	// API docs (no auth)
	mux.HandleFunc("GET /swagger/", h.Docs.SwaggerUI)
	mux.HandleFunc("GET /docs/openapi.yaml", h.Docs.Spec)

	// Public
	mux.HandleFunc("POST /api/v1/auth/login", h.Auth.Login)

	// Wajib Pajak
	mux.Handle("GET /api/v1/wajib-pajak",
		chain(http.HandlerFunc(h.WajibPajak.List), authMW, all, middleware.SelfOnly))
	mux.Handle("POST /api/v1/wajib-pajak",
		chain(http.HandlerFunc(h.WajibPajak.Create), authMW, staff))
	mux.Handle("GET /api/v1/wajib-pajak/{id}",
		chain(http.HandlerFunc(h.WajibPajak.GetByID), authMW, all))
	mux.Handle("PUT /api/v1/wajib-pajak/{id}",
		chain(http.HandlerFunc(h.WajibPajak.Update), authMW, staff))

	// Kewajiban Pajak
	mux.Handle("GET /api/v1/kewajiban-pajak",
		chain(http.HandlerFunc(h.Kewajiban.List), authMW, all, middleware.SelfOnly))
	mux.Handle("POST /api/v1/kewajiban-pajak",
		chain(http.HandlerFunc(h.Kewajiban.Create), authMW, admin))

	// Pembayaran
	mux.Handle("POST /api/v1/pembayaran",
		chain(http.HandlerFunc(h.Pembayaran.Create), authMW, staff))

	// Laporan
	mux.Handle("GET /api/v1/laporan",
		chain(http.HandlerFunc(h.Laporan.Get), authMW, staff))

	// Denda
	mux.Handle("GET /api/v1/denda/{id}",
		chain(http.HandlerFunc(h.Denda.Get), authMW, all))

	return mux
}
