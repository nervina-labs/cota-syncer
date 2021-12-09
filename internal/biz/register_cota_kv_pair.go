package biz

import (
	"context"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/logger"
)

type RegisterCotaKvPair struct {
	BlockNumber uint64
	LockHash    string
	LockHashCRC uint32
}

type RegisterCotaKvPairRepo interface {
	CreateRegisterCotaKvPair(ctx context.Context, register *RegisterCotaKvPair) error
	DeleteRegisterCotaKvPairs(ctx context.Context, blockNumber uint64) error
}

type RegisterCotaKvPairUsecase struct {
	repo   RegisterCotaKvPairRepo
	logger *logger.Logger
}

func NewRegisterCotaKvPairUsecase(repo RegisterCotaKvPairRepo, logger *logger.Logger) *RegisterCotaKvPairUsecase {
	return &RegisterCotaKvPairUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *RegisterCotaKvPairUsecase) Create(ctx context.Context, register *RegisterCotaKvPair) error {
	return uc.repo.CreateRegisterCotaKvPair(ctx, register)
}

func (uc *RegisterCotaKvPairUsecase) DeleteByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteRegisterCotaKvPairs(ctx, blockNumber)
}
