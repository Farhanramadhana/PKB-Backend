package dto

import "errors"

var (
	ErrNotFound            = errors.New("data tidak ditemukan")
	ErrUnauthorized        = errors.New("tidak terautentikasi")
	ErrForbidden           = errors.New("akses ditolak")
	ErrConflict            = errors.New("data sudah ada")
	ErrWajibPajakInaktif   = errors.New("wajib pajak tidak aktif")
	ErrTanggalMasaDepan    = errors.New("tanggal bayar tidak boleh di masa depan")
	ErrJumlahHarusPositif  = errors.New("jumlah bayar harus lebih dari nol")
	ErrInvalidCredentials  = errors.New("username atau password salah")
	ErrInvalidToken        = errors.New("token tidak valid")
)
