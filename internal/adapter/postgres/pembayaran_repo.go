package postgres

import (
	"context"
	"errors"

	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/bpka/tps-pkb/internal/domain/entity"
	"github.com/bpka/tps-pkb/internal/domain/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type pembayaranRepo struct{ pool *pgxpool.Pool }

func NewPembayaranRepository(pool *pgxpool.Pool) port.PembayaranRepository {
	return &pembayaranRepo{pool: pool}
}

// SaveWithDenda atomically writes: pembayaran, denda records, and updates kewajiban status.
func (r *pembayaranRepo) SaveWithDenda(
	ctx context.Context,
	p *entity.Pembayaran,
	dendaList []entity.Denda,
	newTotal decimal.Decimal,
	newStatus entity.StatusKewajiban,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	_, err = tx.Exec(ctx,
		`INSERT INTO pembayaran (id, kewajiban_id, wajib_pajak_id, user_id, tanggal_bayar, jumlah_bayar, status, catatan_petugas)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		p.ID, p.KewajibanID, p.WajibPajakID, p.UserID,
		p.TanggalBayar.Format("2006-01-02"), p.JumlahBayar, p.Status, p.CatatanPetugas)
	if err != nil {
		return err
	}

	for _, d := range dendaList {
		_, err = tx.Exec(ctx,
			`INSERT INTO denda (id, kewajiban_id, pembayaran_id, jenis, dasar, tarif, jumlah)
			 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			d.ID, d.KewajibanID, d.PembayaranID, d.Jenis, d.Dasar, d.Tarif, d.Jumlah)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(ctx,
		`UPDATE kewajiban_pajak SET total_dibayar=$1, status=$2, updated_at=NOW() WHERE id=$3`,
		newTotal, newStatus, p.KewajibanID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pembayaranRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Pembayaran, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, kewajiban_id, wajib_pajak_id, user_id, tanggal_bayar, jumlah_bayar, status, catatan_petugas, created_at
		 FROM pembayaran WHERE id=$1`, id)
	var p entity.Pembayaran
	err := row.Scan(&p.ID, &p.KewajibanID, &p.WajibPajakID, &p.UserID,
		&p.TanggalBayar, &p.JumlahBayar, &p.Status, &p.CatatanPetugas, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dto.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}
