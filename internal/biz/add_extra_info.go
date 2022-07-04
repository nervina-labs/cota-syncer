package biz

import (
	"context"

	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
)

type AddExtraInfoRepo interface {
	CreateWithdrawTxHash(ctx context.Context, txHash string) error
	DeleteWithdrawTxHash(ctx context.Context, blockNumber uint64) error
	CreateDefineLockScriptId(ctx context.Context, lockScriptId uint) error
	DeleteDefineLockScriptId(ctx context.Context, blockNumber uint64) error
	FindOrCreateScript(ctx context.Context, script *Script) error
}

type AddExtraInfoUsecase struct {
	repo   AddExtraInfoRepo
	logger *logger.Logger
}

func NewAddExtraInfoRepo(repo AddExtraInfoRepo, logger *logger.Logger) *AddExtraInfoUsecase {
	return &AddExtraInfoUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *AddExtraInfoUsecase) CreateWithdrawTxHash(ctx context.Context, txHash string) error {
	return uc.repo.CreateWithdrawTxHash(ctx, txHash)
}

func (uc *AddExtraInfoUsecase) DeleteWithdrawTxHash(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteWithdrawTxHash(ctx, blockNumber)
}

func (uc *AddExtraInfoUsecase) CreateDefineLockScriptId(ctx context.Context, lockScriptId uint) error {
	return uc.repo.CreateDefineLockScriptId(ctx, lockScriptId)
}

func (uc *AddExtraInfoUsecase) DeleteDefineLockScriptId(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteDefineLockScriptId(ctx, blockNumber)
}

func (uc *AddExtraInfoUsecase) FindOrCreateScript(ctx context.Context, script *Script) error {
	return uc.repo.FindOrCreateScript(ctx, script)
}
