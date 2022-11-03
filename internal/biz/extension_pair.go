package biz

import (
	"context"
	"time"

	"github.com/nervina-labs/cota-syncer/internal/logger"
)

type ExtensionPair struct {
	BlockNumber uint64
	LockHash    string
	LockHashCRC uint32
	Key         string
	Value       string
	TxIndex     uint32
	UpdatedAt   time.Time
}

type ExtensionPairRepo interface {
	CreateExtensionPair(ctx context.Context, extension *ExtensionPair) error
	DeleteExtensionPairs(ctx context.Context, blockNumber uint64) error
	ParseExtensionPairs(blockNumber uint64, entry Entry) ([]ExtensionPair, []SubKeyPair, error)
}

type ExtensionPairUsecase struct {
	repo   ExtensionPairRepo
	logger *logger.Logger
}

func NewExtensionPairUsecase(repo ExtensionPairRepo, logger *logger.Logger) *ExtensionPairUsecase {
	return &ExtensionPairUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *ExtensionPairUsecase) Create(ctx context.Context, extension *ExtensionPair) error {
	return uc.repo.CreateExtensionPair(ctx, extension)
}

func (uc *ExtensionPairUsecase) DeleteByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteExtensionPairs(ctx, blockNumber)
}

func (uc *ExtensionPairUsecase) ParseExtensionPair(blockNumber uint64, entry Entry) ([]ExtensionPair, []SubKeyPair, error) {
	return uc.repo.ParseExtensionPairs(blockNumber, entry)
}
