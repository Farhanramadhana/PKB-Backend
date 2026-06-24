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

func buildPembayaranService(ctrl *gomock.Controller) (
	*mock.MockKewajibanRepository,
	*mock.MockWajibPajakRepository,
	*mock.MockPembayaranRepository,
	*mock.MockDendaService,
	*mock.MockAuditService,
) {
	return mock.NewMockKewajibanRepository(ctrl),
		mock.NewMockWajibPajakRepository(ctrl),
		mock.NewMockPembayaranRepository(ctrl),
		mock.NewMockDendaService(ctrl),
		mock.NewMockAuditService(ctrl)
}

func TestPembayaranService_Create_JumlahNol(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kewajibanRepo, wpRepo, pembayaranRepo, dendaSvc, auditSvc := buildPembayaranService(ctrl)
	svc := service.NewPembayaranService(kewajibanRepo, wpRepo, pembayaranRepo, dendaSvc, auditSvc)

	_, _, err := svc.Create(t.Context(), dto.CreatePembayaranRequest{
		KewajibanID:  uuid.New(),
		TanggalBayar: time.Now().Format("2006-01-02"),
		JumlahBayar:  "0",
	}, uuid.New())

	assert.ErrorIs(t, err, dto.ErrJumlahHarusPositif)
}

func TestPembayaranService_Create_TanggalMasaDepan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kewajibanRepo, wpRepo, pembayaranRepo, dendaSvc, auditSvc := buildPembayaranService(ctrl)
	svc := service.NewPembayaranService(kewajibanRepo, wpRepo, pembayaranRepo, dendaSvc, auditSvc)

	_, _, err := svc.Create(t.Context(), dto.CreatePembayaranRequest{
		KewajibanID:  uuid.New(),
		TanggalBayar: time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
		JumlahBayar:  "500000",
	}, uuid.New())

	assert.ErrorIs(t, err, dto.ErrTanggalMasaDepan)
}

func TestPembayaranService_Create_WajibPajakInaktif(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kewajibanRepo, wpRepo, pembayaranRepo, dendaSvc, auditSvc := buildPembayaranService(ctrl)

	wpID := uuid.New()
	kewajibanID := uuid.New()
	kewajiban := &entity.KewajibanPajak{
		ID:           kewajibanID,
		WajibPajakID: wpID,
		PokokPajak:   decimal.NewFromInt(1_000_000),
		PeriodeFinal: time.Now().AddDate(0, 1, 0),
		TotalDibayar: decimal.Zero,
	}
	kewajibanRepo.EXPECT().FindByID(gomock.Any(), kewajibanID).Return(kewajiban, nil)
	wpRepo.EXPECT().FindByID(gomock.Any(), wpID).Return(&entity.WajibPajak{IsActive: false}, nil)

	svc := service.NewPembayaranService(kewajibanRepo, wpRepo, pembayaranRepo, dendaSvc, auditSvc)
	_, _, err := svc.Create(t.Context(), dto.CreatePembayaranRequest{
		KewajibanID:  kewajibanID,
		TanggalBayar: time.Now().Format("2006-01-02"),
		JumlahBayar:  "500000",
	}, uuid.New())

	assert.ErrorIs(t, err, dto.ErrWajibPajakInaktif)
}

func TestPembayaranService_Create_DendaGabungan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kewajibanRepo, wpRepo, pembayaranRepo, dendaSvc, auditSvc := buildPembayaranService(ctrl)

	wpID := uuid.New()
	kewajibanID := uuid.New()
	kewajiban := &entity.KewajibanPajak{
		ID:           kewajibanID,
		WajibPajakID: wpID,
		PokokPajak:   decimal.NewFromInt(1_000_000),
		PeriodeFinal: time.Now().AddDate(0, -1, 0), // overdue
		TotalDibayar: decimal.Zero,
	}

	expectedDenda := []entity.Denda{
		{Jenis: entity.DendaTelatBayar, Jumlah: decimal.NewFromInt(20_000)},
		{Jenis: entity.DendaKurangBayar, Jumlah: decimal.NewFromInt(5_000)},
	}

	kewajibanRepo.EXPECT().FindByID(gomock.Any(), kewajibanID).Return(kewajiban, nil)
	wpRepo.EXPECT().FindByID(gomock.Any(), wpID).Return(&entity.WajibPajak{ID: wpID, IsActive: true}, nil)
	dendaSvc.EXPECT().Hitung(gomock.Any()).Return(expectedDenda)
	pembayaranRepo.EXPECT().
		SaveWithDenda(gomock.Any(), gomock.Any(), gomock.Len(2), gomock.Any(), gomock.Any()).
		Return(nil)
	auditSvc.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil)

	svc := service.NewPembayaranService(kewajibanRepo, wpRepo, pembayaranRepo, dendaSvc, auditSvc)
	pembayaran, dendaList, err := svc.Create(t.Context(), dto.CreatePembayaranRequest{
		KewajibanID:  kewajibanID,
		TanggalBayar: time.Now().Format("2006-01-02"),
		JumlahBayar:  "500000",
	}, uuid.New())

	require.NoError(t, err)
	assert.NotNil(t, pembayaran)
	assert.Equal(t, entity.StatusPembayaranKurangBayar, pembayaran.Status)
	assert.Len(t, dendaList, 2)
}

func TestPembayaranService_Create_Lunas(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kewajibanRepo, wpRepo, pembayaranRepo, dendaSvc, auditSvc := buildPembayaranService(ctrl)

	wpID := uuid.New()
	kewajibanID := uuid.New()
	kewajiban := &entity.KewajibanPajak{
		ID:           kewajibanID,
		WajibPajakID: wpID,
		PokokPajak:   decimal.NewFromInt(1_000_000),
		PeriodeFinal: time.Now().AddDate(0, 1, 0), // not overdue
		TotalDibayar: decimal.Zero,
	}

	kewajibanRepo.EXPECT().FindByID(gomock.Any(), kewajibanID).Return(kewajiban, nil)
	wpRepo.EXPECT().FindByID(gomock.Any(), wpID).Return(&entity.WajibPajak{ID: wpID, IsActive: true}, nil)
	dendaSvc.EXPECT().Hitung(gomock.Any()).Return(nil) // no fines
	pembayaranRepo.EXPECT().
		SaveWithDenda(gomock.Any(), gomock.Any(), gomock.Len(0), gomock.Any(), entity.StatusLunas).
		Return(nil)
	auditSvc.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil)

	svc := service.NewPembayaranService(kewajibanRepo, wpRepo, pembayaranRepo, dendaSvc, auditSvc)
	pembayaran, dendaList, err := svc.Create(t.Context(), dto.CreatePembayaranRequest{
		KewajibanID:  kewajibanID,
		TanggalBayar: time.Now().Format("2006-01-02"),
		JumlahBayar:  "1000000",
	}, uuid.New())

	require.NoError(t, err)
	assert.Equal(t, entity.StatusPembayaranLunas, pembayaran.Status)
	assert.Empty(t, dendaList)
}
