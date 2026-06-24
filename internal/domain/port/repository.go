package port

import (
	"context"

	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/bpka/tps-pkb/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

//go:generate mockgen -source=repository.go -destination=../../mock/mock_repository.go -package=mock

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
}

type WajibPajakRepository interface {
	Save(ctx context.Context, wp *entity.WajibPajak) error
	Update(ctx context.Context, wp *entity.WajibPajak) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.WajibPajak, error)
	FindAll(ctx context.Context, filter dto.WajibPajakFilter) ([]entity.WajibPajak, int, error)
	ExistsByNIK(ctx context.Context, nik string) (bool, error)
	ExistsByNPWP(ctx context.Context, npwp string) (bool, error)
}

type KendaraanRepository interface {
	Save(ctx context.Context, k *entity.Kendaraan) error
	FindByNomorPolisi(ctx context.Context, nomorPolisi string) (*entity.Kendaraan, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Kendaraan, error)
}

type KewajibanRepository interface {
	Save(ctx context.Context, k *entity.KewajibanPajak) error
	Update(ctx context.Context, k *entity.KewajibanPajak) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.KewajibanPajak, error)
	FindAll(ctx context.Context, filter dto.KewajibanFilter) ([]entity.KewajibanPajak, int, error)
}

type PembayaranRepository interface {
	// SaveWithDenda persists pembayaran + denda + kewajiban status update atomically.
	SaveWithDenda(ctx context.Context, p *entity.Pembayaran, dendaList []entity.Denda, newTotal decimal.Decimal, newStatus entity.StatusKewajiban) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Pembayaran, error)
}

type DendaRepository interface {
	FindByKewajibanID(ctx context.Context, kewajibanID uuid.UUID) ([]entity.Denda, error)
	FindByPembayaranID(ctx context.Context, pembayaranID uuid.UUID) ([]entity.Denda, error)
}

type LaporanRepository interface {
	GetLaporan(ctx context.Context, filter dto.LaporanFilter) (*entity.Laporan, error)
}

type AuditRepository interface {
	Save(ctx context.Context, entry AuditEntry, userID *uuid.UUID, username, ipAddress, userAgent string) error
}
