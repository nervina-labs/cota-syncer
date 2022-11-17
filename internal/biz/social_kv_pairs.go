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
	CreateSocialKeyPair(ctx context.Context, extension *SocialKvPair) error
	DeleteSocialKeyPairs(ctx context.Context, blockNumber uint64) error
}

type SocialPairRepoUsecase struct {
	repo   SocialPairRepo
	logger *logger.Logger
}

func NewSocialKeyPairRepoUsecase(repo SocialPairRepo, logger *logger.Logger) *SocialPairRepoUsecase {
	return &SocialPairRepoUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *SocialPairRepoUsecase) Create(ctx context.Context, social *SocialKvPair) error {
	return uc.repo.CreateSocialKeyPair(ctx, social)
}

func (uc *SocialPairRepoUsecase) DeleteByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteSocialKeyPairs(ctx, blockNumber)
}
