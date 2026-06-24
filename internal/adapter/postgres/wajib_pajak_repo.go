package postgres

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/bpka/tps-pkb/internal/domain/entity"
	"github.com/bpka/tps-pkb/internal/domain/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type wajibPajakRepo struct{ pool *pgxpool.Pool }

func NewWajibPajakRepository(pool *pgxpool.Pool) port.WajibPajakRepository {
	return &wajibPajakRepo{pool: pool}
}

func (r *wajibPajakRepo) Save(ctx context.Context, wp *entity.WajibPajak) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO wajib_pajak (id, nama, jenis, nik, npwp, nib, alamat, no_telp, email, is_active)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		wp.ID, wp.Nama, wp.Jenis, wp.NIK, wp.NPWP, wp.NIB, wp.Alamat, wp.NoTelp, wp.Email, wp.IsActive)
	return err
}

func (r *wajibPajakRepo) Update(ctx context.Context, wp *entity.WajibPajak) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE wajib_pajak SET nama=$1, alamat=$2, no_telp=$3, email=$4, is_active=$5, updated_at=NOW()
		 WHERE id=$6`,
		wp.Nama, wp.Alamat, wp.NoTelp, wp.Email, wp.IsActive, wp.ID)
	return err
}

func (r *wajibPajakRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.WajibPajak, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, nama, jenis, nik, npwp, nib, alamat, no_telp, email, is_active, created_at, updated_at
		 FROM wajib_pajak WHERE id=$1`, id)
	return scanWajibPajak(row)
}

func (r *wajibPajakRepo) FindAll(ctx context.Context, filter dto.WajibPajakFilter) ([]entity.WajibPajak, int, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	q := psql.Select("id,nama,jenis,nik,npwp,nib,alamat,no_telp,email,is_active,created_at,updated_at").
		From("wajib_pajak")
	cq := psql.Select("COUNT(*)").From("wajib_pajak")

	if filter.Nama != "" {
		like := fmt.Sprintf("%%%s%%", filter.Nama)
		q = q.Where(sq.ILike{"nama": like})
		cq = cq.Where(sq.ILike{"nama": like})
	}
	if filter.Jenis != "" {
		q = q.Where(sq.Eq{"jenis": filter.Jenis})
		cq = cq.Where(sq.Eq{"jenis": filter.Jenis})
	}
	if filter.IsActive != nil {
		q = q.Where(sq.Eq{"is_active": *filter.IsActive})
		cq = cq.Where(sq.Eq{"is_active": *filter.IsActive})
	}
	if filter.OwnerID != nil {
		q = q.Where(sq.Eq{"id": *filter.OwnerID})
		cq = cq.Where(sq.Eq{"id": *filter.OwnerID})
	}

	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit
	q = q.OrderBy("created_at DESC").Limit(uint64(filter.Limit)).Offset(uint64(offset))

	// count
	csql, cargs, err := cq.ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int
	if err := r.pool.QueryRow(ctx, csql, cargs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// data
	qsql, qargs, err := q.ToSql()
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, qsql, qargs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []entity.WajibPajak
	for rows.Next() {
		var wp entity.WajibPajak
		if err := rows.Scan(&wp.ID, &wp.Nama, &wp.Jenis, &wp.NIK, &wp.NPWP, &wp.NIB,
			&wp.Alamat, &wp.NoTelp, &wp.Email, &wp.IsActive, &wp.CreatedAt, &wp.UpdatedAt); err != nil {
			return nil, 0, err
		}
		result = append(result, wp)
	}
	return result, total, nil
}

func (r *wajibPajakRepo) ExistsByNIK(ctx context.Context, nik string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM wajib_pajak WHERE nik=$1)`, nik).Scan(&exists)
	return exists, err
}

func (r *wajibPajakRepo) ExistsByNPWP(ctx context.Context, npwp string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM wajib_pajak WHERE npwp=$1)`, npwp).Scan(&exists)
	return exists, err
}

func scanWajibPajak(row pgx.Row) (*entity.WajibPajak, error) {
	var wp entity.WajibPajak
	err := row.Scan(&wp.ID, &wp.Nama, &wp.Jenis, &wp.NIK, &wp.NPWP, &wp.NIB,
		&wp.Alamat, &wp.NoTelp, &wp.Email, &wp.IsActive, &wp.CreatedAt, &wp.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dto.ErrNotFound
		}
		return nil, err
	}
	return &wp, nil
}
