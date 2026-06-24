package service

import (
	"context"

	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/bpka/tps-pkb/internal/domain/entity"
	"github.com/bpka/tps-pkb/internal/domain/port"
)

type laporanService struct {
	repo port.LaporanRepository
}

func NewLaporanService(repo port.LaporanRepository) port.LaporanService {
	return &laporanService{repo: repo}
}

func (s *laporanService) Get(ctx context.Context, filter dto.LaporanFilter) (*entity.Laporan, error) {
	return s.repo.GetLaporan(ctx, filter)
}
