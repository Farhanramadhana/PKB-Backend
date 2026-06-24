package postgres

import (
	"context"

	"github.com/bpka/tps-pkb/internal/domain/entity"
	"github.com/bpka/tps-pkb/internal/domain/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type dendaRepo struct{ pool *pgxpool.Pool }

func NewDendaRepository(pool *pgxpool.Pool) port.DendaRepository {
	return &dendaRepo{pool: pool}
}

func (r *dendaRepo) FindByKewajibanID(ctx context.Context, kewajibanID uuid.UUID) ([]entity.Denda, error) {
	return r.query(ctx, `SELECT id, kewajiban_id, pembayaran_id, jenis, dasar, tarif, jumlah, created_at
		 FROM denda WHERE kewajiban_id=$1 ORDER BY created_at DESC`, kewajibanID)
}

func (r *dendaRepo) FindByPembayaranID(ctx context.Context, pembayaranID uuid.UUID) ([]entity.Denda, error) {
	return r.query(ctx, `SELECT id, kewajiban_id, pembayaran_id, jenis, dasar, tarif, jumlah, created_at
		 FROM denda WHERE pembayaran_id=$1 ORDER BY created_at DESC`, pembayaranID)
}

func (r *dendaRepo) query(ctx context.Context, sql string, arg uuid.UUID) ([]entity.Denda, error) {
	rows, err := r.pool.Query(ctx, sql, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.Denda
	for rows.Next() {
		var d entity.Denda
		if err := rows.Scan(&d.ID, &d.KewajibanID, &d.PembayaranID, &d.Jenis, &d.Dasar, &d.Tarif, &d.Jumlah, &d.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}
