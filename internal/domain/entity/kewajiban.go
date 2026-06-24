package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type StatusKewajiban string

const (
	StatusBelumBayar  StatusKewajiban = "BELUM_BAYAR"
	StatusLunas       StatusKewajiban = "LUNAS"
	StatusKurangBayar StatusKewajiban = "KURANG_BAYAR"
	StatusLebihBayar  StatusKewajiban = "LEBIH_BAYAR"
)

type KewajibanPajak struct {
	ID           uuid.UUID
	KendaraanID  uuid.UUID
	WajibPajakID uuid.UUID
	TahunPajak   int
	PeriodeAwal  time.Time
	PeriodeFinal time.Time
	PokokPajak   decimal.Decimal
	Status       StatusKewajiban
	TotalDibayar decimal.Decimal
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Joined fields (populated on queries that JOIN)
	NomorPolisi    string
	WajibPajakNama string
}
