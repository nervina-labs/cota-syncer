package biz

import (
	"context"

	"github.com/nervina-labs/cota-syncer/internal/logger"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
)

type RegisterCotaKvPair struct {
	BlockNumber uint64
	LockHash    string
	CotaCellID  uint64
	LockScriptId uint
}

type RegisterCotaKvPairRepo interface {
	CreateRegisterCotaKvPair(ctx context.Context, register *RegisterCotaKvPair) error
	DeleteRegisterCotaKvPairs(ctx context.Context, blockNumber uint64) error
	ParseRegistryEntries(ctx context.Context, blockNumber uint64, tx *ckbTypes.Transaction) ([]RegisterCotaKvPair, error)
	FindOrCreateScript(ctx context.Context, script *Script) error
}

type RegisterCotaKvPairUsecase struct {
	repo   RegisterCotaKvPairRepo
	logger *logger.Logger
}

func NewRegisterCotaKvPairUsecase(repo RegisterCotaKvPairRepo, logger *logger.Logger) *RegisterCotaKvPairUsecase {
	return &RegisterCotaKvPairUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *RegisterCotaKvPairUsecase) Create(ctx context.Context, register *RegisterCotaKvPair) error {
	return uc.repo.CreateRegisterCotaKvPair(ctx, register)
}

func (uc *RegisterCotaKvPairUsecase) DeleteByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteRegisterCotaKvPairs(ctx, blockNumber)
}

func (uc *RegisterCotaKvPairUsecase) ParseRegistryEntries(ctx context.Context, blockNumber uint64, tx *ckbTypes.Transaction) ([]RegisterCotaKvPair, error) {
	return uc.repo.ParseRegistryEntries(ctx, blockNumber, tx)
}

func (uc *RegisterCotaKvPairUsecase) FindOrCreateScript(ctx context.Context, script *Script) error {
	return uc.repo.FindOrCreateScript(ctx, script)
}
