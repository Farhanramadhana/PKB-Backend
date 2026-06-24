package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type JenisDenda string

const (
	DendaTelatBayar JenisDenda = "TELAT_BAYAR"
	DendaKurangBayar JenisDenda = "KURANG_BAYAR"
)

type Denda struct {
	ID           uuid.UUID
	KewajibanID  uuid.UUID
	PembayaranID uuid.UUID
	Jenis        JenisDenda
	Dasar        decimal.Decimal
	Tarif        decimal.Decimal
	Jumlah       decimal.Decimal
	CreatedAt    time.Time
}
