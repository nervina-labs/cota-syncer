package biz

import (
	"context"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	"time"
)

type DefineCotaNftKvPair struct {
	BlockNumber uint64
	CotaId      string
	Total       uint32
	Issued      uint32
	Configure   uint8
	LockHash    string
	LockHashCRC uint32
	TxIndex     uint32
	UpdatedAt   time.Time
}

type DefineCotaNftKvPairRepo interface {
	CreateDefineCotaNftKvPair(ctx context.Context, d *DefineCotaNftKvPair) error
	DeleteDefineCotaNftKvPairs(ctx context.Context, blockNumber uint64) error
	ParseDefineCotaEntries(blockNumber uint64, entry Entry) ([]DefineCotaNftKvPair, error)
}

type DefineCotaNftKvPairUsecase struct {
	repo   DefineCotaNftKvPairRepo
	logger *logger.Logger
}

func NewDefineCotaNftKvPairUsecase(repo DefineCotaNftKvPairRepo, logger *logger.Logger) *DefineCotaNftKvPairUsecase {
	return &DefineCotaNftKvPairUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *DefineCotaNftKvPairUsecase) Create(ctx context.Context, d *DefineCotaNftKvPair) error {
	return uc.repo.CreateDefineCotaNftKvPair(ctx, d)
}

func (uc *DefineCotaNftKvPairUsecase) DeleteByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteDefineCotaNftKvPairs(ctx, blockNumber)
}

func (uc *DefineCotaNftKvPairUsecase) ParseDefineCotaEntries(blockNumber uint64, entry Entry) ([]DefineCotaNftKvPair, error) {
	return uc.repo.ParseDefineCotaEntries(blockNumber, entry)
}
