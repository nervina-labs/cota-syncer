package biz

import (
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
)

type MintCotaKvPair struct{}

type MintCotaKvPairRepo interface {
	ParseMintCotaEntries(blockNumber uint64, entry Entry) ([]DefineCotaNftKvPair, []WithdrawCotaNftKvPair, error)
}

type MintCotaKvPairUsecase struct {
	repo   MintCotaKvPairRepo
	logger *logger.Logger
}

func NewMintCotaKvPairUsecase(repo MintCotaKvPairRepo, logger *logger.Logger) *MintCotaKvPairUsecase {
	return &MintCotaKvPairUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *MintCotaKvPairUsecase) ParseMintCotaEntries(blockNumber uint64, entry Entry) ([]DefineCotaNftKvPair, []WithdrawCotaNftKvPair, error) {
	return uc.repo.ParseMintCotaEntries(blockNumber, entry)
}
