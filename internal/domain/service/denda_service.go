package service

import (
	"context"
	"time"

	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/bpka/tps-pkb/internal/domain/entity"
	"github.com/bpka/tps-pkb/internal/domain/port"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	tarifTelatBayar  = decimal.NewFromFloat(0.02)
	tarifKurangBayar = decimal.NewFromFloat(0.01)
)

type dendaService struct {
	repo port.DendaRepository
}

func NewDendaService(repo port.DendaRepository) port.DendaService {
	return &dendaService{repo: repo}
}

// Hitung is a pure function — no I/O, no context, deterministic.
// It applies the three denda rules from the spec:
//   - Denda Telat Bayar: 2% × pokok_pajak when tanggal_bayar > periode_final
//   - Denda Kurang Bayar: 1% × selisih when total_dibayar + jumlah_bayar < pokok_pajak
//   - Denda Gabungan: both when late AND underpaid
func (s *dendaService) Hitung(input dto.DendaInput) []entity.Denda {
	var result []entity.Denda

	tanggal := input.TanggalBayar.Truncate(24 * time.Hour)
	jatuhTempo := input.PeriodeFinal.Truncate(24 * time.Hour)

	if tanggal.After(jatuhTempo) {
		jumlah := input.PokokPajak.Mul(tarifTelatBayar)
		result = append(result, entity.Denda{
			ID:           uuid.New(),
			KewajibanID:  input.KewajibanID,
			PembayaranID: input.PembayaranID,
			Jenis:        entity.DendaTelatBayar,
			Dasar:        input.PokokPajak,
			Tarif:        tarifTelatBayar,
			Jumlah:       jumlah,
			CreatedAt:    time.Now(),
		})
	}

	newTotal := input.TotalDibayar.Add(input.JumlahBayar)
	if newTotal.LessThan(input.PokokPajak) {
		selisih := input.PokokPajak.Sub(newTotal)
		jumlah := selisih.Mul(tarifKurangBayar)
		result = append(result, entity.Denda{
			ID:           uuid.New(),
			KewajibanID:  input.KewajibanID,
			PembayaranID: input.PembayaranID,
			Jenis:        entity.DendaKurangBayar,
			Dasar:        selisih,
			Tarif:        tarifKurangBayar,
			Jumlah:       jumlah,
			CreatedAt:    time.Now(),
		})
	}

	return result
}

func (s *dendaService) GetByKewajibanID(ctx context.Context, id uuid.UUID) ([]entity.Denda, error) {
	return s.repo.FindByKewajibanID(ctx, id)
}
