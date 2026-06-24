package dto

import "github.com/google/uuid"

type CreateKewajibanRequest struct {
	WajibPajakID uuid.UUID `json:"wajib_pajak_id" validate:"required"`
	NomorPolisi  string    `json:"nomor_polisi"   validate:"required"`
	Merk         string    `json:"merk"           validate:"required"`
	Model        string    `json:"model"          validate:"required"`
	Tahun        int       `json:"tahun"          validate:"required,min=1900"`
	TahunPajak   int       `json:"tahun_pajak"    validate:"required,min=2000"`
	PeriodeAwal  string    `json:"periode_awal"   validate:"required"`
	PeriodeFinal string    `json:"periode_final"  validate:"required"`
	PokokPajak   string    `json:"pokok_pajak"    validate:"required"`
}

type KewajibanFilter struct {
	WajibPajakID *uuid.UUID
	Status       string
	TahunPajak   *int
	Page         int
	Limit        int
}
