package biz

import (
	"context"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	"time"
)

type SubKeyPair struct {
	BlockNumber uint64
	LockHash    string
	SubType     string
	ExtData     uint32
	AlgIndex    uint16
	PubkeyHash  string
	UpdatedAt   time.Time
}

type SubKeyPairRepo interface {
	CreateSubKeyPair(ctx context.Context, extension *SubKeyPair) error
	DeleteSubKeyPairs(ctx context.Context, blockNumber uint64) error
	ParseSubKeyPairs(blockNumber uint64, entry Entry) ([]SubKeyPair, error)
}

type SubKeyPairRepoUsecase struct {
	repo   SubKeyPairRepo
	logger *logger.Logger
}

func NewSubKeyPairRepoUsecase(repo SubKeyPairRepo, logger *logger.Logger) *SubKeyPairRepoUsecase {
	return &SubKeyPairRepoUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *SubKeyPairRepoUsecase) Create(ctx context.Context, extension *SubKeyPair) error {
	return uc.repo.CreateSubKeyPair(ctx, extension)
}

func (uc *SubKeyPairRepoUsecase) DeleteByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteSubKeyPairs(ctx, blockNumber)
}

func (uc *SubKeyPairRepoUsecase) ParseExtensionPair(blockNumber uint64, entry Entry) ([]SubKeyPair, error) {
	return uc.repo.ParseSubKeyPairs(blockNumber, entry)
}
