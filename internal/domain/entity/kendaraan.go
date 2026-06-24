package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Kendaraan struct {
	ID             uuid.UUID
	WajibPajakID   uuid.UUID
	NomorPolisi    string
	Merk           string
	Model          string
	Tahun          int
	JenisKendaraan string
	BPKB           string
	STNK           string
	NilaiJual      decimal.Decimal
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
