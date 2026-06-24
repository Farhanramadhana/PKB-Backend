package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type StatusPembayaran string

const (
	StatusPembayaranLunas      StatusPembayaran = "LUNAS"
	StatusPembayaranKurangBayar StatusPembayaran = "KURANG_BAYAR"
	StatusPembayaranLebihBayar  StatusPembayaran = "LEBIH_BAYAR"
)

type Pembayaran struct {
	ID              uuid.UUID
	KewajibanID     uuid.UUID
	WajibPajakID    uuid.UUID
	UserID          uuid.UUID
	TanggalBayar    time.Time
	JumlahBayar     decimal.Decimal
	Status          StatusPembayaran
	CatatanPetugas  string
	CreatedAt       time.Time
}
