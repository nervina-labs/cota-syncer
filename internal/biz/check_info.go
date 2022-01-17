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
	Id          uint64
	BlockNumber uint64
	BlockHash   string
	CheckType   CheckType
}

type CheckInfoRepo interface {
	FindLastCheckInfo(ctx context.Context, info *CheckInfo) error
	CreateCheckInfo(ctx context.Context, info *CheckInfo) error
	CleanCheckInfo(ctx context.Context) error
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

func (uc *CheckInfoUsecase) Clean(ctx context.Context) error {
	return uc.repo.CleanCheckInfo(ctx)
}
