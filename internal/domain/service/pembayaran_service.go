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

type pembayaranService struct {
	kewajibanRepo port.KewajibanRepository
	wajibPajakRepo port.WajibPajakRepository
	pembayaranRepo port.PembayaranRepository
	dendaSvc      port.DendaService
	audit         port.AuditService
}

func NewPembayaranService(
	kewajibanRepo port.KewajibanRepository,
	wajibPajakRepo port.WajibPajakRepository,
	pembayaranRepo port.PembayaranRepository,
	dendaSvc port.DendaService,
	audit port.AuditService,
) port.PembayaranService {
	return &pembayaranService{
		kewajibanRepo:  kewajibanRepo,
		wajibPajakRepo: wajibPajakRepo,
		pembayaranRepo: pembayaranRepo,
		dendaSvc:       dendaSvc,
		audit:          audit,
	}
}

func (s *pembayaranService) Create(ctx context.Context, req dto.CreatePembayaranRequest, userID uuid.UUID) (*entity.Pembayaran, []entity.Denda, error) {
	jumlah, err := decimal.NewFromString(req.JumlahBayar)
	if err != nil || jumlah.LessThanOrEqual(decimal.Zero) {
		return nil, nil, dto.ErrJumlahHarusPositif
	}

	tanggal, err := time.Parse("2006-01-02", req.TanggalBayar)
	if err != nil {
		return nil, nil, dto.ErrTanggalMasaDepan
	}
	if tanggal.After(time.Now().Truncate(24 * time.Hour)) {
		return nil, nil, dto.ErrTanggalMasaDepan
	}

	kewajiban, err := s.kewajibanRepo.FindByID(ctx, req.KewajibanID)
	if err != nil {
		return nil, nil, err
	}

	wp, err := s.wajibPajakRepo.FindByID(ctx, kewajiban.WajibPajakID)
	if err != nil {
		return nil, nil, err
	}
	if !wp.IsActive {
		return nil, nil, dto.ErrWajibPajakInaktif
	}

	newTotal := kewajiban.TotalDibayar.Add(jumlah)
	status := pembayaranStatusFromAmounts(newTotal, kewajiban.PokokPajak)

	pembayaran := &entity.Pembayaran{
		ID:             uuid.New(),
		KewajibanID:    kewajiban.ID,
		WajibPajakID:   kewajiban.WajibPajakID,
		UserID:         userID,
		TanggalBayar:   tanggal,
		JumlahBayar:    jumlah,
		Status:         status,
		CatatanPetugas: req.CatatanPetugas,
		CreatedAt:      time.Now(),
	}

	dendaList := s.dendaSvc.Hitung(dto.DendaInput{
		KewajibanID:  kewajiban.ID,
		PembayaranID: pembayaran.ID,
		PokokPajak:   kewajiban.PokokPajak,
		PeriodeFinal: kewajiban.PeriodeFinal,
		TanggalBayar: tanggal,
		JumlahBayar:  jumlah,
		TotalDibayar: kewajiban.TotalDibayar,
	})

	newKewajibanStatus := kewajibanStatusFromAmounts(newTotal, kewajiban.PokokPajak)

	if err := s.pembayaranRepo.SaveWithDenda(ctx, pembayaran, dendaList, newTotal, newKewajibanStatus); err != nil {
		return nil, nil, err
	}

	_ = s.audit.Log(ctx, port.AuditEntry{
		TableName: "pembayaran",
		RecordID:  pembayaran.ID,
		Action:    "CREATE",
		NewData:   pembayaran,
	})

	return pembayaran, dendaList, nil
}

func pembayaranStatusFromAmounts(dibayar, pokok decimal.Decimal) entity.StatusPembayaran {
	switch {
	case dibayar.Equal(pokok):
		return entity.StatusPembayaranLunas
	case dibayar.LessThan(pokok):
		return entity.StatusPembayaranKurangBayar
	default:
		return entity.StatusPembayaranLebihBayar
	}
}

func kewajibanStatusFromAmounts(dibayar, pokok decimal.Decimal) entity.StatusKewajiban {
	switch {
	case dibayar.Equal(pokok):
		return entity.StatusLunas
	case dibayar.LessThan(pokok):
		return entity.StatusKurangBayar
	default:
		return entity.StatusLebihBayar
	}
}
