package biz

import (
	"context"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/logger"
)

type HoldCotaNftKVPair struct {
	BlockNumber    uint64
	CotaId         string
	CotaIdCRC      uint32
	Total          uint32
	TokenIndex     uint32
	State          uint8
	Configure      uint8
	Characteristic string
	LockHash       string
	LockHashCRC    uint32
}

type HoldCotaNftKVPairRepo interface {
	CreateHoldCotaNftKVPair(ctx context.Context, h *HoldCotaNftKVPair) error
	DeleteHoldCotaNftKVPairs(ctx context.Context, blockNumber uint64) error
}

type HoldCotaNftKVPairUsecase struct {
	repo   HoldCotaNftKVPairRepo
	logger *logger.Logger
}

func NewHoldCotaNftKVPairUsecase(repo HoldCotaNftKVPairRepo, logger *logger.Logger) *HoldCotaNftKVPairUsecase {
	return &HoldCotaNftKVPairUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *HoldCotaNftKVPairUsecase) Create(ctx context.Context, h *HoldCotaNftKVPair) error {
	return uc.repo.CreateHoldCotaNftKVPair(ctx, h)
}

func (uc *HoldCotaNftKVPairUsecase) DeleteByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteHoldCotaNftKVPairs(ctx, blockNumber)
}
