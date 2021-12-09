package biz

import (
	"context"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/logger"
)

type HoldCotaNftKvPair struct {
	BlockNumber    uint64
	CotaId         string
	CotaIdCRC      uint32
	TokenIndex     uint32
	State          uint8
	Configure      uint8
	Characteristic string
	LockHash       string
	LockHashCRC    uint32
}

type HoldCotaNftKvPairRepo interface {
	CreateHoldCotaNftKvPair(ctx context.Context, h *HoldCotaNftKvPair) error
	DeleteHoldCotaNftKvPairs(ctx context.Context, blockNumber uint64) error
}

type HoldCotaNftKvPairUsecase struct {
	repo   HoldCotaNftKvPairRepo
	logger *logger.Logger
}

func NewHoldCotaNftKvPairUsecase(repo HoldCotaNftKvPairRepo, logger *logger.Logger) *HoldCotaNftKvPairUsecase {
	return &HoldCotaNftKvPairUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *HoldCotaNftKvPairUsecase) Create(ctx context.Context, h *HoldCotaNftKvPair) error {
	return uc.repo.CreateHoldCotaNftKvPair(ctx, h)
}

func (uc *HoldCotaNftKvPairUsecase) DeleteByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteHoldCotaNftKvPairs(ctx, blockNumber)
}
