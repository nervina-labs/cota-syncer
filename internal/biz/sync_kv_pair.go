package biz

import (
	"context"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/logger"
)

type KvPair struct {
	Registers     []RegisterCotaKvPair
	DefineCotas   []DefineCotaNftKvPair
	HoldCotas     []HoldCotaNftKvPair
	WithdrawCotas []WithdrawCotaNftKvPair
	ClaimedCotas  []ClaimedCotaNftKvPair
}

type KvPairRepo interface {
	CreateKvPairs(ctx context.Context, kvPair *KvPair) error
	DeleteKvPairs(ctx context.Context, blockNumber uint64) error
}

type SyncKvPairUsecase struct {
	repo   KvPairRepo
	logger *logger.Logger
}

func NewSyncKvPairUsecase(repo KvPairRepo, logger *logger.Logger) *SyncKvPairUsecase {
	return &SyncKvPairUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc SyncKvPairUsecase) CreateKvPairs(ctx context.Context, kvPair *KvPair) error {
	return uc.repo.CreateKvPairs(ctx, kvPair)
}

func (uc SyncKvPairUsecase) DeleteKvPairs(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteKvPairs(ctx, blockNumber)
}
