package biz

import (
	"context"
	"github.com/nervina-labs/cota-syncer/internal/logger"
)

type CheckType uint8

const (
	SyncBlock    CheckType = iota // SyncBlock = 0
	SyncMetadata                  // SyncMetadata = 1
)

func (t CheckType) String() string {
	return []string{"sync_block_event", "sync_metadata_event"}[t]
}

type CheckInfo struct {
	Id          uint64
	BlockNumber uint64
	BlockHash   string
	CheckType   CheckType
}

type CheckInfoRepo interface {
	FindLastCheckInfo(ctx context.Context, info *CheckInfo) error
	CreateCheckInfo(ctx context.Context, info *CheckInfo) error
	CleanCheckInfo(ctx context.Context, checkType CheckType) error
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

func (uc *CheckInfoUsecase) LastCheckInfo(ctx context.Context, checkInfo *CheckInfo) error {
	return uc.repo.FindLastCheckInfo(ctx, checkInfo)
}

func (uc *CheckInfoUsecase) Create(ctx context.Context, checkInfo *CheckInfo) error {
	return uc.repo.CreateCheckInfo(ctx, checkInfo)
}

func (uc *CheckInfoUsecase) Clean(ctx context.Context, checkType CheckType) error {
	return uc.repo.CleanCheckInfo(ctx, checkType)
}
