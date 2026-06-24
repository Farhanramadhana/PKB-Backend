package dto

import "github.com/google/uuid"

type CreateWajibPajakRequest struct {
	Nama   string `json:"nama"   validate:"required,min=3,max=255"`
	Jenis  string `json:"jenis"  validate:"required,oneof=INDIVIDU BADAN_USAHA"`
	NIK    string `json:"nik"    validate:"required_if=Jenis INDIVIDU,omitempty,nik"`
	NPWP   string `json:"npwp"   validate:"required_if=Jenis BADAN_USAHA,omitempty,npwp"`
	NIB    string `json:"nib"    validate:"required_if=Jenis BADAN_USAHA,omitempty,min=1"`
	Alamat string `json:"alamat" validate:"required"`
	NoTelp string `json:"no_telp"`
	Email  string `json:"email"  validate:"omitempty,email"`
}

type UpdateWajibPajakRequest struct {
	Nama   string `json:"nama"   validate:"omitempty,min=3,max=255"`
	Alamat string `json:"alamat"`
	NoTelp string `json:"no_telp"`
	Email  string `json:"email"  validate:"omitempty,email"`
}

type WajibPajakFilter struct {
	Nama     string
	Jenis    string
	IsActive *bool
	Page     int
	Limit    int
	OwnerID  *uuid.UUID // non-nil when role is WAJIB_PAJAK
}
