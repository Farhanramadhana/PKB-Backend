package port

import (
	"context"

	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/bpka/tps-pkb/internal/domain/entity"
	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=../../mock/mock_service.go -package=mock

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.TokenResponse, error)
}

type WajibPajakService interface {
	Create(ctx context.Context, req dto.CreateWajibPajakRequest) (*entity.WajibPajak, error)
	List(ctx context.Context, filter dto.WajibPajakFilter) ([]entity.WajibPajak, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.WajibPajak, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateWajibPajakRequest) (*entity.WajibPajak, error)
	SetActive(ctx context.Context, id uuid.UUID, active bool) error
}

type KewajibanService interface {
	Create(ctx context.Context, req dto.CreateKewajibanRequest) (*entity.KewajibanPajak, error)
	List(ctx context.Context, filter dto.KewajibanFilter) ([]entity.KewajibanPajak, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.KewajibanPajak, error)
}

type PembayaranService interface {
	Create(ctx context.Context, req dto.CreatePembayaranRequest, userID uuid.UUID) (*entity.Pembayaran, []entity.Denda, error)
}

type DendaService interface {
	// Hitung is a pure calculation — no I/O. Safe to call without context.
	Hitung(input dto.DendaInput) []entity.Denda
	GetByKewajibanID(ctx context.Context, id uuid.UUID) ([]entity.Denda, error)
}

type LaporanService interface {
	Get(ctx context.Context, filter dto.LaporanFilter) (*entity.Laporan, error)
}
