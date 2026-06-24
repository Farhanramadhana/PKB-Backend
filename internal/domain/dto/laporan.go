package dto

import "github.com/google/uuid"

type LaporanFilter struct {
	StartDate    string
	EndDate      string
	Status       string
	WajibPajakID *uuid.UUID
	Page         int
	Limit        int
}
