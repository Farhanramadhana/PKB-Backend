package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	adapthttp "github.com/bpka/tps-pkb/internal/adapter/http"
	"github.com/bpka/tps-pkb/internal/adapter/http/handler"
	"github.com/bpka/tps-pkb/internal/adapter/postgres"
	"github.com/bpka/tps-pkb/internal/domain/service"
	"github.com/bpka/tps-pkb/internal/infrastructure/config"
	"github.com/bpka/tps-pkb/internal/infrastructure/database"
	"github.com/bpka/tps-pkb/internal/infrastructure/token"
	"github.com/go-playground/validator/v10"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Database
	pool, err := database.NewPool(ctx, cfg.Database.URL)
	if err != nil {
		slog.Error("database connection failed", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Migrations
	if err := database.RunMigrations(cfg.Database.URL, "migrations"); err != nil {
		slog.Error("migrations failed", "err", err)
		os.Exit(1)
	}
	slog.Info("migrations applied")

	// Infrastructure
	tokenProvider := token.NewJWTProvider(cfg.JWT.Secret, cfg.JWT.TTLHours)
	validate := newValidator()

	// Repositories
	userRepo       := postgres.NewUserRepository(pool)
	wajibPajakRepo := postgres.NewWajibPajakRepository(pool)
	kendaraanRepo  := postgres.NewKendaraanRepository(pool)
	kewajibanRepo  := postgres.NewKewajibanRepository(pool)
	pembayaranRepo := postgres.NewPembayaranRepository(pool)
	dendaRepo      := postgres.NewDendaRepository(pool)
	laporanRepo    := postgres.NewLaporanRepository(pool)
	auditRepo      := postgres.NewAuditRepository(pool)

	// Audit service
	auditSvc := postgres.NewAuditService(auditRepo)

	// Domain services
	authSvc       := service.NewAuthService(userRepo, tokenProvider)
	wajibPajakSvc := service.NewWajibPajakService(wajibPajakRepo, auditSvc)
	dendaSvc      := service.NewDendaService(dendaRepo)
	kewajibanSvc  := service.NewKewajibanService(kewajibanRepo, kendaraanRepo, wajibPajakRepo, auditSvc)
	pembayaranSvc := service.NewPembayaranService(kewajibanRepo, wajibPajakRepo, pembayaranRepo, dendaSvc, auditSvc)
	laporanSvc    := service.NewLaporanService(laporanRepo)

	// Handlers
	handlers := &adapthttp.Handlers{
		Auth:       handler.NewAuthHandler(authSvc, validate),
		WajibPajak: handler.NewWajibPajakHandler(wajibPajakSvc, validate),
		Kewajiban:  handler.NewKewajibanHandler(kewajibanSvc, validate),
		Pembayaran: handler.NewPembayaranHandler(pembayaranSvc, validate),
		Denda:      handler.NewDendaHandler(dendaSvc),
		Laporan:    handler.NewLaporanHandler(laporanSvc),
		Health:     handler.NewHealthHandler(),
		Docs:       handler.NewDocsHandler(),
	}

	router := adapthttp.NewRouter(handlers, tokenProvider)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "port", cfg.Server.Port, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("forced shutdown", "err", err)
	}
	slog.Info("server stopped")
}

func newValidator() *validator.Validate {
	v := validator.New()

	// NIK: exactly 16 numeric digits
	_ = v.RegisterValidation("nik", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		if len(s) != 16 {
			return false
		}
		for _, c := range s {
			if c < '0' || c > '9' {
				return false
			}
		}
		return true
	})

	// NPWP: exactly 15 numeric digits
	_ = v.RegisterValidation("npwp", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		if len(s) != 15 {
			return false
		}
		for _, c := range s {
			if c < '0' || c > '9' {
				return false
			}
		}
		return true
	})

	return v
}
