package biz

import (
	"context"
	"github.com/nervina-labs/cota-syncer/internal/logger"
)

type InvalidDataRepo interface {
	Clean(ctx context.Context, blockNumber uint64) error
}

type InvalidDataUsecase struct {
	repo   InvalidDataRepo
	logger *logger.Logger
}

func NewInvalidDataUsecase(repo InvalidDataRepo, logger *logger.Logger) *InvalidDataUsecase {
	return &InvalidDataUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *InvalidDataUsecase) Clean(ctx context.Context, blockNumber uint64) error {
	return uc.repo.Clean(ctx, blockNumber)
}
