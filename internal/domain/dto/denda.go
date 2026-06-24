package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type DendaInput struct {
	KewajibanID  uuid.UUID
	PembayaranID uuid.UUID
	PokokPajak   decimal.Decimal
	PeriodeFinal time.Time
	TanggalBayar time.Time
	JumlahBayar  decimal.Decimal
	TotalDibayar decimal.Decimal
}
