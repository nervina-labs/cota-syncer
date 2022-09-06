package biz

import (
	"context"

	"github.com/nervina-labs/cota-syncer/internal/logger"
)

type TransferCotaKvPair struct{}

type TransferCotaKvPairRepo interface {
	ParseTransferCotaEntries(blockNumber uint64, entry Entry) ([]ClaimedCotaNftKvPair, []WithdrawCotaNftKvPair, error)
	ParseTransferUpdateCotaEntries(blockNumber uint64, entry Entry) ([]ClaimedCotaNftKvPair, []WithdrawCotaNftKvPair, error)
	FindOrCreateScript(ctx context.Context, script *Script) error
}

type TransferCotaKvPairUsecase struct {
	repo   TransferCotaKvPairRepo
	logger *logger.Logger
}

func NewTransferCotaKvPairUsecase(repo TransferCotaKvPairRepo, logger *logger.Logger) *TransferCotaKvPairUsecase {
	return &TransferCotaKvPairUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *TransferCotaKvPairUsecase) ParseTransferCotaEntries(blockNumber uint64, entry Entry) ([]ClaimedCotaNftKvPair, []WithdrawCotaNftKvPair, error) {
	return uc.repo.ParseTransferCotaEntries(blockNumber, entry)
}

func (uc *TransferCotaKvPairUsecase) ParseTransferUpdateCotaEntries(blockNumber uint64, entry Entry) ([]ClaimedCotaNftKvPair, []WithdrawCotaNftKvPair, error) {
	return uc.repo.ParseTransferUpdateCotaEntries(blockNumber, entry)
}

func (uc *TransferCotaKvPairUsecase) FindOrCreateScript(ctx context.Context, script *Script) error {
	return uc.repo.FindOrCreateScript(ctx, script)
}
