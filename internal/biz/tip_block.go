package biz

import (
	"context"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/logger"
)

type TipBlock struct {
	BlockNumber uint64
	BlockHash   string
}

type TipBlockRepo interface {
	FindOrCreateTipBlock(ctx context.Context, tipBlock *TipBlock) error
	UpdateTipBlock(ctx context.Context, tipBlock *TipBlock) error
}

type TipBlockUsecase struct {
	repo   TipBlockRepo
	logger *logger.Logger
}

func NewTipBlockUsecase(repo TipBlockRepo, logger *logger.Logger) *TipBlockUsecase {
	return &TipBlockUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *TipBlockUsecase) FindOrCreate(ctx context.Context, tipBlock *TipBlock) error {
	return uc.repo.FindOrCreateTipBlock(ctx, tipBlock)
}

func (uc *TipBlockUsecase) Update(ctx context.Context, tipBlock *TipBlock) error {
	return uc.Update(ctx, tipBlock)
}
