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

type kewajibanService struct {
	repo        port.KewajibanRepository
	kendaraanRepo port.KendaraanRepository
	wajibPajakRepo port.WajibPajakRepository
	audit       port.AuditService
}

func NewKewajibanService(
	repo port.KewajibanRepository,
	kendaraanRepo port.KendaraanRepository,
	wajibPajakRepo port.WajibPajakRepository,
	audit port.AuditService,
) port.KewajibanService {
	return &kewajibanService{
		repo:           repo,
		kendaraanRepo:  kendaraanRepo,
		wajibPajakRepo: wajibPajakRepo,
		audit:          audit,
	}
}

func (s *kewajibanService) Create(ctx context.Context, req dto.CreateKewajibanRequest) (*entity.KewajibanPajak, error) {
	wp, err := s.wajibPajakRepo.FindByID(ctx, req.WajibPajakID)
	if err != nil {
		return nil, err
	}
	if !wp.IsActive {
		return nil, dto.ErrWajibPajakInaktif
	}

	kendaraan, err := s.kendaraanRepo.FindByNomorPolisi(ctx, req.NomorPolisi)
	if err != nil {
		pokok, parseErr := decimal.NewFromString(req.PokokPajak)
		if parseErr != nil || pokok.LessThanOrEqual(decimal.Zero) {
			return nil, dto.ErrJumlahHarusPositif
		}
		kendaraan = &entity.Kendaraan{
			ID:           uuid.New(),
			WajibPajakID: req.WajibPajakID,
			NomorPolisi:  req.NomorPolisi,
			Merk:         req.Merk,
			Model:        req.Model,
			Tahun:        req.Tahun,
			NilaiJual:    pokok,
		}
		if saveErr := s.kendaraanRepo.Save(ctx, kendaraan); saveErr != nil {
			return nil, saveErr
		}
	}

	periodeAwal, err := time.Parse("2006-01-02", req.PeriodeAwal)
	if err != nil {
		return nil, err
	}
	periodeFinal, err := time.Parse("2006-01-02", req.PeriodeFinal)
	if err != nil {
		return nil, err
	}
	pokok, err := decimal.NewFromString(req.PokokPajak)
	if err != nil || pokok.LessThanOrEqual(decimal.Zero) {
		return nil, dto.ErrJumlahHarusPositif
	}

	kewajiban := &entity.KewajibanPajak{
		ID:           uuid.New(),
		KendaraanID:  kendaraan.ID,
		WajibPajakID: req.WajibPajakID,
		TahunPajak:   req.TahunPajak,
		PeriodeAwal:  periodeAwal,
		PeriodeFinal: periodeFinal,
		PokokPajak:   pokok,
		Status:       entity.StatusBelumBayar,
		TotalDibayar: decimal.Zero,
	}

	if err := s.repo.Save(ctx, kewajiban); err != nil {
		return nil, err
	}

	_ = s.audit.Log(ctx, port.AuditEntry{
		TableName: "kewajiban_pajak",
		RecordID:  kewajiban.ID,
		Action:    "CREATE",
		NewData:   kewajiban,
	})

	return kewajiban, nil
}

func (s *kewajibanService) List(ctx context.Context, filter dto.KewajibanFilter) ([]entity.KewajibanPajak, int, error) {
	return s.repo.FindAll(ctx, filter)
}

func (s *kewajibanService) GetByID(ctx context.Context, id uuid.UUID) (*entity.KewajibanPajak, error) {
	return s.repo.FindByID(ctx, id)
}
