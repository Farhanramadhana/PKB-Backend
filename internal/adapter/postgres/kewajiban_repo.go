package postgres

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/bpka/tps-pkb/internal/domain/entity"
	"github.com/bpka/tps-pkb/internal/domain/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type kewajibanRepo struct{ pool *pgxpool.Pool }

func NewKewajibanRepository(pool *pgxpool.Pool) port.KewajibanRepository {
	return &kewajibanRepo{pool: pool}
}

func (r *kewajibanRepo) Save(ctx context.Context, k *entity.KewajibanPajak) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO kewajiban_pajak
		 (id, kendaraan_id, wajib_pajak_id, tahun_pajak, periode_awal, periode_final, pokok_pajak, status, total_dibayar)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		k.ID, k.KendaraanID, k.WajibPajakID, k.TahunPajak,
		k.PeriodeAwal.Format("2006-01-02"), k.PeriodeFinal.Format("2006-01-02"),
		k.PokokPajak, k.Status, k.TotalDibayar)
	return err
}

func (r *kewajibanRepo) Update(ctx context.Context, k *entity.KewajibanPajak) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE kewajiban_pajak SET status=$1, total_dibayar=$2, updated_at=NOW() WHERE id=$3`,
		k.Status, k.TotalDibayar, k.ID)
	return err
}

func (r *kewajibanRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.KewajibanPajak, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT kp.id, kp.kendaraan_id, kp.wajib_pajak_id, kp.tahun_pajak,
		        kp.periode_awal, kp.periode_final, kp.pokok_pajak, kp.status, kp.total_dibayar,
		        kp.created_at, kp.updated_at,
		        k.nomor_polisi, wp.nama
		 FROM kewajiban_pajak kp
		 JOIN kendaraan k ON k.id = kp.kendaraan_id
		 JOIN wajib_pajak wp ON wp.id = kp.wajib_pajak_id
		 WHERE kp.id=$1`, id)
	return scanKewajiban(row)
}

func (r *kewajibanRepo) FindAll(ctx context.Context, filter dto.KewajibanFilter) ([]entity.KewajibanPajak, int, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	base := `kewajiban_pajak kp
		JOIN kendaraan k ON k.id = kp.kendaraan_id
		JOIN wajib_pajak wp ON wp.id = kp.wajib_pajak_id`

	cols := "kp.id,kp.kendaraan_id,kp.wajib_pajak_id,kp.tahun_pajak,kp.periode_awal,kp.periode_final,kp.pokok_pajak,kp.status,kp.total_dibayar,kp.created_at,kp.updated_at,k.nomor_polisi,wp.nama"
	q := psql.Select(cols).From(base)
	cq := psql.Select("COUNT(*)").From(base)

	if filter.WajibPajakID != nil {
		q = q.Where(sq.Eq{"kp.wajib_pajak_id": *filter.WajibPajakID})
		cq = cq.Where(sq.Eq{"kp.wajib_pajak_id": *filter.WajibPajakID})
	}
	if filter.Status != "" {
		q = q.Where(sq.Eq{"kp.status": filter.Status})
		cq = cq.Where(sq.Eq{"kp.status": filter.Status})
	}
	if filter.TahunPajak != nil {
		q = q.Where(sq.Eq{"kp.tahun_pajak": *filter.TahunPajak})
		cq = cq.Where(sq.Eq{"kp.tahun_pajak": *filter.TahunPajak})
	}

	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit
	q = q.OrderBy("kp.created_at DESC").Limit(uint64(filter.Limit)).Offset(uint64(offset))

	csql, cargs, err := cq.ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int
	if err := r.pool.QueryRow(ctx, csql, cargs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	qsql, qargs, err := q.ToSql()
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, qsql, qargs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []entity.KewajibanPajak
	for rows.Next() {
		k, err := scanKewajibanRow(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, *k)
	}
	return result, total, nil
}

type kewajibanScanner interface {
	Scan(dest ...any) error
}

func scanKewajiban(row pgx.Row) (*entity.KewajibanPajak, error) {
	k, err := scanKewajibanRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dto.ErrNotFound
		}
		return nil, err
	}
	return k, nil
}

func scanKewajibanRow(scanner kewajibanScanner) (*entity.KewajibanPajak, error) {
	var k entity.KewajibanPajak
	err := scanner.Scan(
		&k.ID, &k.KendaraanID, &k.WajibPajakID, &k.TahunPajak,
		&k.PeriodeAwal, &k.PeriodeFinal, &k.PokokPajak, &k.Status, &k.TotalDibayar,
		&k.CreatedAt, &k.UpdatedAt, &k.NomorPolisi, &k.WajibPajakNama,
	)
	return &k, err
}
