package postgres

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/bpka/tps-pkb/internal/adapter/http/middleware"
	"github.com/bpka/tps-pkb/internal/domain/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type auditRepo struct{ pool *pgxpool.Pool }

func NewAuditRepository(pool *pgxpool.Pool) port.AuditRepository {
	return &auditRepo{pool: pool}
}

type pgAuditService struct {
	repo port.AuditRepository
}

// NewAuditService wraps the repository to implement port.AuditService.
// It extracts user identity from context and delegates persistence to AuditRepository.
func NewAuditService(repo port.AuditRepository) port.AuditService {
	return &pgAuditService{repo: repo}
}

func (s *pgAuditService) Log(ctx context.Context, entry port.AuditEntry) error {
	claims := middleware.ClaimsFromContext(ctx)
	var userID *uuid.UUID
	username := "system"
	if claims != nil {
		userID = &claims.UserID
		username = claims.Username
	}
	ip := middleware.IPFromContext(ctx)
	ua := middleware.UserAgentFromContext(ctx)

	if err := s.repo.Save(ctx, entry, userID, username, ip, ua); err != nil {
		// Audit failure must not break the main operation — log and continue.
		slog.Error("audit log failed", "table", entry.TableName, "action", entry.Action, "err", err)
	}
	return nil
}

func (r *auditRepo) Save(ctx context.Context, entry port.AuditEntry, userID *uuid.UUID, username, ipAddress, userAgent string) error {
	oldData, _ := json.Marshal(entry.OldData)
	newData, _ := json.Marshal(entry.NewData)

	var oldJSON, newJSON *[]byte
	if entry.OldData != nil {
		oldJSON = &oldData
	}
	if entry.NewData != nil {
		newJSON = &newData
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO audit_logs (table_name, record_id, action, user_id, username, old_data, new_data, ip_address, user_agent)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		entry.TableName, entry.RecordID, entry.Action, userID, username, oldJSON, newJSON, ipAddress, userAgent)
	return err
}
