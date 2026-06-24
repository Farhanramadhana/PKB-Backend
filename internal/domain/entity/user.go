package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin      Role = "ADMIN"
	RolePetugas    Role = "PETUGAS"
	RoleWajibPajak Role = "WAJIB_PAJAK"
)

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
	Role         Role
	WajibPajakID *uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
