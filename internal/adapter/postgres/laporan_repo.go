package postgres

import (
	"context"
	"fmt"

	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/bpka/tps-pkb/internal/domain/entity"
	"github.com/bpka/tps-pkb/internal/domain/port"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type laporanRepo struct{ pool *pgxpool.Pool }

func NewLaporanRepository(pool *pgxpool.Pool) port.LaporanRepository {
	return &laporanRepo{pool: pool}
}

func (r *laporanRepo) GetLaporan(ctx context.Context, filter dto.LaporanFilter) (*entity.Laporan, error) {
	args := []any{}
	where := "WHERE 1=1"
	argIdx := 1

	if filter.StartDate != "" {
		where += fmt.Sprintf(" AND p.tanggal_bayar >= $%d", argIdx)
		args = append(args, filter.StartDate)
		argIdx++
	}
	if filter.EndDate != "" {
		where += fmt.Sprintf(" AND p.tanggal_bayar <= $%d", argIdx)
		args = append(args, filter.EndDate)
		argIdx++
	}
	if filter.Status != "" {
		where += fmt.Sprintf(" AND kp.status = $%d", argIdx)
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.WajibPajakID != nil {
		where += fmt.Sprintf(" AND kp.wajib_pajak_id = $%d", argIdx)
		args = append(args, *filter.WajibPajakID)
		argIdx++
	}

	summarySQL := fmt.Sprintf(`
		SELECT
			COALESCE(SUM(kp.pokok_pajak), 0)    AS total_kewajiban,
			COALESCE(SUM(kp.total_dibayar), 0)  AS total_dibayar,
			COALESCE(SUM(d.total_denda), 0)     AS total_denda
		FROM kewajiban_pajak kp
		LEFT JOIN pembayaran p ON p.kewajiban_id = kp.id
		LEFT JOIN (
			SELECT kewajiban_id, SUM(jumlah) AS total_denda
			FROM denda GROUP BY kewajiban_id
		) d ON d.kewajiban_id = kp.id
		%s`, where)

	var totalKewajiban, totalDibayar, totalDenda decimal.Decimal
	if err := r.pool.QueryRow(ctx, summarySQL, args...).Scan(&totalKewajiban, &totalDibayar, &totalDenda); err != nil {
		return nil, err
	}

	itemsSQL := fmt.Sprintf(`
		SELECT
			kp.id, wp.nama, k.nomor_polisi, kp.tahun_pajak,
			kp.pokok_pajak, kp.total_dibayar,
			COALESCE(d.total_denda, 0) AS total_denda,
			kp.status
		FROM kewajiban_pajak kp
		JOIN wajib_pajak wp ON wp.id = kp.wajib_pajak_id
		JOIN kendaraan k ON k.id = kp.kendaraan_id
		LEFT JOIN pembayaran p ON p.kewajiban_id = kp.id
		LEFT JOIN (
			SELECT kewajiban_id, SUM(jumlah) AS total_denda
			FROM denda GROUP BY kewajiban_id
		) d ON d.kewajiban_id = kp.id
		%s
		GROUP BY kp.id, wp.nama, k.nomor_polisi, d.total_denda
		ORDER BY kp.created_at DESC`, where)

	rows, err := r.pool.Query(ctx, itemsSQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.LaporanItem
	for rows.Next() {
		var item entity.LaporanItem
		if err := rows.Scan(&item.KewajibanID, &item.WajibPajakNama, &item.NomorPolisi,
			&item.TahunPajak, &item.PokokPajak, &item.TotalDibayar, &item.TotalDenda, &item.Status); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return &entity.Laporan{
		TotalKewajiban: totalKewajiban,
		TotalDibayar:   totalDibayar,
		TotalDenda:     totalDenda,
		SisaKewajiban:  totalKewajiban.Sub(totalDibayar),
		Items:          items,
	}, nil
}
