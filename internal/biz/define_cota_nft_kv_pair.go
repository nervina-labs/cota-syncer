package biz

import (
	"context"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/logger"
)

type DefineCotaNftKVPair struct {
	BlockNumber uint64
	CotaId      string
	CotaIdCRC   uint32
	Total       uint32
	Issued      uint32
	Configure   uint8
	LockHash    string
	LockHashCRC uint32
}

type DefineCotaNftKVPairRepo interface {
	CreateDefineCotaNftKVPair(ctx context.Context, d *DefineCotaNftKVPair) error
	DeleteDefineCotaNftKVPairs(ctx context.Context, blockNumber uint64) error
}

type DefineCotaNftKVPairUsecase struct {
	repo   DefineCotaNftKVPairRepo
	logger *logger.Logger
}

func NewDefineCotaNftKVPairUsecase(repo DefineCotaNftKVPairRepo, logger *logger.Logger) *DefineCotaNftKVPairUsecase {
	return &DefineCotaNftKVPairUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *DefineCotaNftKVPairUsecase) Create(ctx context.Context, d *DefineCotaNftKVPair) error {
	return uc.repo.CreateDefineCotaNftKVPair(ctx, d)
}

func (uc *DefineCotaNftKVPairUsecase) DeleteByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteDefineCotaNftKVPairs(ctx, blockNumber)
}
