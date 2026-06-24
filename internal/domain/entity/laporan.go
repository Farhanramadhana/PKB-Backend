package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Laporan struct {
	TotalKewajiban decimal.Decimal
	TotalDibayar   decimal.Decimal
	TotalDenda     decimal.Decimal
	SisaKewajiban  decimal.Decimal
	Items          []LaporanItem
}

type LaporanItem struct {
	KewajibanID    uuid.UUID
	WajibPajakNama string
	NomorPolisi    string
	TahunPajak     int
	PokokPajak     decimal.Decimal
	TotalDibayar   decimal.Decimal
	TotalDenda     decimal.Decimal
	Status         StatusKewajiban
}
