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
)

type kendaraanRepo struct{ pool *pgxpool.Pool }

func NewKendaraanRepository(pool *pgxpool.Pool) port.KendaraanRepository {
	return &kendaraanRepo{pool: pool}
}

func (r *kendaraanRepo) Save(ctx context.Context, k *entity.Kendaraan) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO kendaraan (id, wajib_pajak_id, nomor_polisi, merk, model, tahun, jenis_kendaraan, bpkb, stnk, nilai_jual)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		k.ID, k.WajibPajakID, k.NomorPolisi, k.Merk, k.Model, k.Tahun,
		k.JenisKendaraan, k.BPKB, k.STNK, k.NilaiJual)
	return err
}

func (r *kendaraanRepo) FindByNomorPolisi(ctx context.Context, nomorPolisi string) (*entity.Kendaraan, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, wajib_pajak_id, nomor_polisi, merk, model, tahun, jenis_kendaraan, bpkb, stnk, nilai_jual, created_at, updated_at
		 FROM kendaraan WHERE nomor_polisi=$1`, nomorPolisi)
	return scanKendaraan(row)
}

func (r *kendaraanRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Kendaraan, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, wajib_pajak_id, nomor_polisi, merk, model, tahun, jenis_kendaraan, bpkb, stnk, nilai_jual, created_at, updated_at
		 FROM kendaraan WHERE id=$1`, id)
	return scanKendaraan(row)
}

func scanKendaraan(row pgx.Row) (*entity.Kendaraan, error) {
	var k entity.Kendaraan
	err := row.Scan(&k.ID, &k.WajibPajakID, &k.NomorPolisi, &k.Merk, &k.Model, &k.Tahun,
		&k.JenisKendaraan, &k.BPKB, &k.STNK, &k.NilaiJual, &k.CreatedAt, &k.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dto.ErrNotFound
		}
		return nil, err
	}
	return &k, nil
}
