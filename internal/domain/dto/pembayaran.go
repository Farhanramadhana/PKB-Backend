package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreatePembayaranRequest struct {
	KewajibanID    uuid.UUID       `json:"kewajiban_id"    validate:"required"`
	TanggalBayar   string          `json:"tanggal_bayar"   validate:"required"`
	JumlahBayar    string          `json:"jumlah_bayar"    validate:"required"`
	CatatanPetugas string          `json:"catatan_petugas"`

	// populated by service after parsing
	ParsedTanggal time.Time       `json:"-"`
	ParsedJumlah  decimal.Decimal `json:"-"`
}

type PembayaranFilter struct {
	WajibPajakID *uuid.UUID
	Page         int
	Limit        int
}
