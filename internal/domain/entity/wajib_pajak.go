package entity

import (
	"time"

	"github.com/google/uuid"
)

type JenisWajibPajak string

const (
	JenisIndividu   JenisWajibPajak = "INDIVIDU"
	JenisBadanUsaha JenisWajibPajak = "BADAN_USAHA"
)

type WajibPajak struct {
	ID        uuid.UUID
	Nama      string
	Jenis     JenisWajibPajak
	NIK       *string
	NPWP      *string
	NIB       *string
	Alamat    string
	NoTelp    string
	Email     string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
