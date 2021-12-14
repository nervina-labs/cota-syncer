package biz

import (
	"context"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
)

type CheckType uint8

const (
	SyncEvent CheckType = iota // SyncEvent = 0
)

func (t CheckType) String() string {
	return []string{"sync_event"}[t]
}

type CheckInfo struct {
	BlockNumber uint64
	BlockHash   string
	CheckType   CheckType
}

type CheckInfoRepo interface {
	FindOrCreateCheckInfo(ctx context.Context, info *CheckInfo) error
	UpdateCheckInfo(ctx context.Context, info CheckInfo) error
}

type CheckInfoUsecase struct {
	repo   CheckInfoRepo
	logger *logger.Logger
}

func NewCheckInfoUsecase(repo CheckInfoRepo, logger *logger.Logger) *CheckInfoUsecase {
	return &CheckInfoUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *CheckInfoUsecase) FindOrCreate(ctx context.Context, checkInfo *CheckInfo) error {
	return uc.repo.FindOrCreateCheckInfo(ctx, checkInfo)
}

func (uc *CheckInfoUsecase) Update(ctx context.Context, checkInfo *CheckInfo) error {
	return uc.Update(ctx, checkInfo)
}
