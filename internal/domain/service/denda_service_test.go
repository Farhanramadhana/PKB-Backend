package service_test

import (
	"testing"
	"time"

	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/bpka/tps-pkb/internal/domain/entity"
	"github.com/bpka/tps-pkb/internal/domain/service"
	"github.com/bpka/tps-pkb/internal/mock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDendaService_Hitung_LunasTepakWaktu(t *testing.T) {
	svc := service.NewDendaService(nil) // Hitung is pure — no repo needed

	result := svc.Hitung(dto.DendaInput{
		KewajibanID:  uuid.New(),
		PembayaranID: uuid.New(),
		PokokPajak:   decimal.NewFromInt(1_000_000),
		PeriodeFinal: time.Now().AddDate(0, 1, 0), // not yet due
		TanggalBayar: time.Now(),
		JumlahBayar:  decimal.NewFromInt(1_000_000),
		TotalDibayar: decimal.Zero,
	})

	assert.Empty(t, result, "no fines when paid in full and on time")
}

func TestDendaService_Hitung_DendaTelatBayar(t *testing.T) {
	svc := service.NewDendaService(nil)

	result := svc.Hitung(dto.DendaInput{
		KewajibanID:  uuid.New(),
		PembayaranID: uuid.New(),
		PokokPajak:   decimal.NewFromInt(1_000_000),
		PeriodeFinal: time.Now().AddDate(0, -1, 0), // 1 month overdue
		TanggalBayar: time.Now(),
		JumlahBayar:  decimal.NewFromInt(1_000_000),
		TotalDibayar: decimal.Zero,
	})

	require.Len(t, result, 1)
	assert.Equal(t, entity.DendaTelatBayar, result[0].Jenis)
	assert.True(t, decimal.NewFromInt(20_000).Equal(result[0].Jumlah), "expected 20000 got %s", result[0].Jumlah)
}

func TestDendaService_Hitung_DendaKurangBayar(t *testing.T) {
	svc := service.NewDendaService(nil)

	result := svc.Hitung(dto.DendaInput{
		KewajibanID:  uuid.New(),
		PembayaranID: uuid.New(),
		PokokPajak:   decimal.NewFromInt(1_000_000),
		PeriodeFinal: time.Now().AddDate(0, 1, 0), // not yet due
		TanggalBayar: time.Now(),
		JumlahBayar:  decimal.NewFromInt(500_000), // underpaid by 500_000
		TotalDibayar: decimal.Zero,
	})

	require.Len(t, result, 1)
	assert.Equal(t, entity.DendaKurangBayar, result[0].Jenis)
	assert.True(t, decimal.NewFromInt(5_000).Equal(result[0].Jumlah), "expected 5000 got %s", result[0].Jumlah)
}

func TestDendaService_Hitung_DendaGabungan(t *testing.T) {
	svc := service.NewDendaService(nil)

	result := svc.Hitung(dto.DendaInput{
		KewajibanID:  uuid.New(),
		PembayaranID: uuid.New(),
		PokokPajak:   decimal.NewFromInt(1_000_000),
		PeriodeFinal: time.Now().AddDate(0, -1, 0), // overdue
		TanggalBayar: time.Now(),
		JumlahBayar:  decimal.NewFromInt(500_000), // also underpaid
		TotalDibayar: decimal.Zero,
	})

	require.Len(t, result, 2, "Denda Gabungan: both fines apply")

	findByJenis := func(jenis entity.JenisDenda) *entity.Denda {
		for i := range result {
			if result[i].Jenis == jenis {
				return &result[i]
			}
		}
		return nil
	}

	telat := findByJenis(entity.DendaTelatBayar)
	require.NotNil(t, telat)
	assert.True(t, decimal.NewFromInt(20_000).Equal(telat.Jumlah), "expected 20000 got %s", telat.Jumlah)

	kurang := findByJenis(entity.DendaKurangBayar)
	require.NotNil(t, kurang)
	assert.True(t, decimal.NewFromInt(5_000).Equal(kurang.Jumlah), "expected 5000 got %s", kurang.Jumlah)
}

func TestDendaService_Hitung_LebihBayar_TidakAdaDenda(t *testing.T) {
	svc := service.NewDendaService(nil)

	// Overpaid on time: no fine of any kind
	result := svc.Hitung(dto.DendaInput{
		KewajibanID:  uuid.New(),
		PembayaranID: uuid.New(),
		PokokPajak:   decimal.NewFromInt(1_000_000),
		PeriodeFinal: time.Now().AddDate(0, 1, 0),
		TanggalBayar: time.Now(),
		JumlahBayar:  decimal.NewFromInt(1_200_000), // overpaid
		TotalDibayar: decimal.Zero,
	})

	assert.Empty(t, result)
}

func TestDendaService_GetByKewajibanID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockDendaRepository(ctrl)
	kewajibanID := uuid.New()
	expected := []entity.Denda{{ID: uuid.New(), KewajibanID: kewajibanID}}

	mockRepo.EXPECT().
		FindByKewajibanID(gomock.Any(), kewajibanID).
		Return(expected, nil)

	svc := service.NewDendaService(mockRepo)
	result, err := svc.GetByKewajibanID(t.Context(), kewajibanID)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}
