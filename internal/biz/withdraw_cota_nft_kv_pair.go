package biz

import (
	"context"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
)

type WithdrawCotaNftKvPair struct {
	BlockNumber         uint64
	CotaId              string
	CotaIdCRC           uint32
	TokenIndex          uint32
	OutPoint            string
	OutPointCrc         uint32
	State               uint8
	Configure           uint8
	Characteristic      string
	ReceiverLockHash    string
	ReceiverLockHashCrc uint32
	LockHash            string
	LockHashCrc         uint32
}

type WithdrawCotaNftKvPairRepo interface {
	CreateWithdrawCotaNftKvPair(ctx context.Context, w *WithdrawCotaNftKvPair) error
	DeleteWithdrawCotaNftKvPairs(ctx context.Context, blockNumber uint64) error
	ParseWithdrawCotaEntries(blockNumber uint64, entry Entry) ([]WithdrawCotaNftKvPair, error)
}

type WithdrawCotaNftKvPairUsecase struct {
	repo   WithdrawCotaNftKvPairRepo
	logger *logger.Logger
}

func NewWithdrawCotaNftKvPairUsecase(repo WithdrawCotaNftKvPairRepo, logger *logger.Logger) *WithdrawCotaNftKvPairUsecase {
	return &WithdrawCotaNftKvPairUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *WithdrawCotaNftKvPairUsecase) Create(ctx context.Context, w *WithdrawCotaNftKvPair) error {
	return uc.repo.CreateWithdrawCotaNftKvPair(ctx, w)
}

func (uc *WithdrawCotaNftKvPairUsecase) DeleteByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteWithdrawCotaNftKvPairs(ctx, blockNumber)
}

func (uc *WithdrawCotaNftKvPairUsecase) ParseWithdrawCotaEntries(blockNumber uint64, entry Entry) ([]WithdrawCotaNftKvPair, error) {
	return uc.repo.ParseWithdrawCotaEntries(blockNumber, entry)
}
