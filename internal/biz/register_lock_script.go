package biz

import (
	"context"

	"github.com/nervina-labs/cota-syncer/internal/logger"
)

type RegisterQueryInfo struct {
	BlockNumber uint64
	LockHash    string
}

type RegisterLockScriptRepo interface {
	AddRegisterLock(ctx context.Context, lockHash string, lockScriptId uint) error
	IsAllHaveLock(ctx context.Context) (bool, error)
	FindRegisterQueryInfos(ctx context.Context, page int, pageSize int) ([]RegisterQueryInfo, error)
	FindOrCreateScript(ctx context.Context, script *Script) error
}

type RegisterLockScriptUsecase struct {
	repo   RegisterLockScriptRepo
	logger *logger.Logger
}

func NewRegisterLockScriptUsecase(repo RegisterLockScriptRepo, logger *logger.Logger) *RegisterLockScriptUsecase {
	return &RegisterLockScriptUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *RegisterLockScriptUsecase) AddRegisterLock(ctx context.Context, lockHash string, lockScriptId uint) error {
	return uc.repo.AddRegisterLock(ctx, lockHash, lockScriptId)
}

func (uc *RegisterLockScriptUsecase) IsAllHaveLock(ctx context.Context) (bool, error) {
	return uc.repo.IsAllHaveLock(ctx)
}

func (uc *RegisterLockScriptUsecase) FindRegisterQueryInfos(ctx context.Context, page int, pageSize int) ([]RegisterQueryInfo, error) {
	return uc.repo.FindRegisterQueryInfos(ctx, page, pageSize)
}

func (uc *RegisterLockScriptUsecase) FindOrCreateScript(ctx context.Context, script *Script) error {
	return uc.repo.FindOrCreateScript(ctx, script)
}
