package service

import (
	"context"

	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/bpka/tps-pkb/internal/domain/entity"
	"github.com/bpka/tps-pkb/internal/domain/port"
	"github.com/google/uuid"
)

type wajibPajakService struct {
	repo  port.WajibPajakRepository
	audit port.AuditService
}

func NewWajibPajakService(repo port.WajibPajakRepository, audit port.AuditService) port.WajibPajakService {
	return &wajibPajakService{repo: repo, audit: audit}
}

func (s *wajibPajakService) Create(ctx context.Context, req dto.CreateWajibPajakRequest) (*entity.WajibPajak, error) {
	if req.Jenis == string(entity.JenisIndividu) && req.NIK != "" {
		exists, err := s.repo.ExistsByNIK(ctx, req.NIK)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, dto.ErrConflict
		}
	}
	if req.Jenis == string(entity.JenisBadanUsaha) && req.NPWP != "" {
		exists, err := s.repo.ExistsByNPWP(ctx, req.NPWP)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, dto.ErrConflict
		}
	}

	jenis := entity.JenisWajibPajak(req.Jenis)
	wp := &entity.WajibPajak{
		ID:       uuid.New(),
		Nama:     req.Nama,
		Jenis:    jenis,
		Alamat:   req.Alamat,
		NoTelp:   req.NoTelp,
		Email:    req.Email,
		IsActive: true,
	}
	if req.NIK != "" {
		wp.NIK = &req.NIK
	}
	if req.NPWP != "" {
		wp.NPWP = &req.NPWP
	}
	if req.NIB != "" {
		wp.NIB = &req.NIB
	}

	if err := s.repo.Save(ctx, wp); err != nil {
		return nil, err
	}

	_ = s.audit.Log(ctx, port.AuditEntry{
		TableName: "wajib_pajak",
		RecordID:  wp.ID,
		Action:    "CREATE",
		NewData:   wp,
	})

	return wp, nil
}

func (s *wajibPajakService) List(ctx context.Context, filter dto.WajibPajakFilter) ([]entity.WajibPajak, int, error) {
	return s.repo.FindAll(ctx, filter)
}

func (s *wajibPajakService) GetByID(ctx context.Context, id uuid.UUID) (*entity.WajibPajak, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *wajibPajakService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateWajibPajakRequest) (*entity.WajibPajak, error) {
	wp, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	old := *wp
	if req.Nama != "" {
		wp.Nama = req.Nama
	}
	if req.Alamat != "" {
		wp.Alamat = req.Alamat
	}
	if req.NoTelp != "" {
		wp.NoTelp = req.NoTelp
	}
	if req.Email != "" {
		wp.Email = req.Email
	}

	if err := s.repo.Update(ctx, wp); err != nil {
		return nil, err
	}

	_ = s.audit.Log(ctx, port.AuditEntry{
		TableName: "wajib_pajak",
		RecordID:  wp.ID,
		Action:    "UPDATE",
		OldData:   old,
		NewData:   wp,
	})

	return wp, nil
}

func (s *wajibPajakService) SetActive(ctx context.Context, id uuid.UUID, active bool) error {
	wp, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	old := *wp
	wp.IsActive = active
	if err := s.repo.Update(ctx, wp); err != nil {
		return err
	}

	_ = s.audit.Log(ctx, port.AuditEntry{
		TableName: "wajib_pajak",
		RecordID:  wp.ID,
		Action:    "UPDATE",
		OldData:   old,
		NewData:   wp,
	})

	return nil
}
