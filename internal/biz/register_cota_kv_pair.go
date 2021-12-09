package biz

import (
	"context"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/logger"
)

type RegisterCotaKVPair struct {
	BlockNumber uint64
	LockHash    string
	LockHashCRC uint32
}

type RegisterCotaKVPairRepo interface {
	CreateRegisterCotaKVPair(ctx context.Context, register *RegisterCotaKVPair) error
	DeleteRegisterCotaKVPairs(ctx context.Context, blockNumber uint64) error
}

type RegisterCotaKVPairUsecase struct {
	repo   RegisterCotaKVPairRepo
	logger *logger.Logger
}

func NewRegisterCotaKVPairUsecase(repo RegisterCotaKVPairRepo, logger *logger.Logger) *RegisterCotaKVPairUsecase {
	return &RegisterCotaKVPairUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *RegisterCotaKVPairUsecase) Create(ctx context.Context, register *RegisterCotaKVPair) error {
	return uc.repo.CreateRegisterCotaKVPair(ctx, register)
}

func (uc *RegisterCotaKVPairUsecase) DeleteByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteRegisterCotaKVPairs(ctx, blockNumber)
}
