package port

import (
	"context"

	"github.com/google/uuid"
)

type AuditEntry struct {
	TableName string
	RecordID  uuid.UUID
	Action    string // CREATE | UPDATE | DELETE
	OldData   any
	NewData   any
}

//go:generate mockgen -source=audit.go -destination=../../mock/mock_audit.go -package=mock
type AuditService interface {
	Log(ctx context.Context, entry AuditEntry) error
}
