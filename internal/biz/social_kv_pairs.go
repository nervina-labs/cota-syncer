package biz

import (
	"context"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	"time"
)

type SocialKvPair struct {
	ID           uint `gorm:"primaryKey"`
	BlockNumber  uint64
	LockHash     string
	LockHashCRC  uint32
	RecoveryMode uint8
	Must         uint8
	Total        uint8
	Signers      string
	UpdatedAt    time.Time
}

type SocialPairRepo interface {
	CreateSocialPair(ctx context.Context, extension *SocialKvPair) error
	DeleteSocialPairs(ctx context.Context, blockNumber uint64) error
}

type SocialPairRepoUsecase struct {
	repo   SocialPairRepo
	logger *logger.Logger
}

func NewSocialPairRepoUsecase(repo SocialPairRepo, logger *logger.Logger) *SocialPairRepoUsecase {
	return &SocialPairRepoUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *SocialPairRepoUsecase) Create(ctx context.Context, social *SocialKvPair) error {
	return uc.repo.CreateSocialPair(ctx, social)
}

func (uc *SocialPairRepoUsecase) DeleteByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteSocialPairs(ctx, blockNumber)
}
