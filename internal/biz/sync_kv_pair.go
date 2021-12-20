package biz

import (
	"context"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
)

type KvPair struct {
	Registers          []RegisterCotaKvPair
	DefineCotas        []DefineCotaNftKvPair
	UpdatedDefineCotas []DefineCotaNftKvPair
	HoldCotas          []HoldCotaNftKvPair
	UpdatedHoldCotas   []HoldCotaNftKvPair
	WithdrawCotas      []WithdrawCotaNftKvPair
	ClaimedCotas       []ClaimedCotaNftKvPair
}

type KvPairRepo interface {
	CreateKvPairs(ctx context.Context, txIndex int, kvPair *KvPair) error
	RestoreKvPairs(ctx context.Context, blockNumber uint64) error
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

func (uc SyncKvPairUsecase) CreateKvPairs(ctx context.Context, txIndex int, kvPair *KvPair) error {
	return uc.repo.CreateKvPairs(ctx, txIndex, kvPair)
}

func (uc SyncKvPairUsecase) RestoreKvPairs(ctx context.Context, blockNumber uint64) error {
	return uc.repo.RestoreKvPairs(ctx, blockNumber)
}
