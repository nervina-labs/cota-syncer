package biz

import (
	"context"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
)

type ClassInfo struct {
	BlockNumber  uint64
	CotaId       string
	Version      string
	Name         string
	Symbol       string
	Description  string
	Image        string
	Audio        string
	Video        string
	Model        string
	Schema       string
	Properties   string
	Localization string
}

type ClassInfoRepo interface {
	CreateClassInfo(ctx context.Context, class *ClassInfo) error
	DeleteClassInfo(ctx context.Context, blockNumber uint64) error
	ParseClassInfo(blockNumber uint64, classJson ClassInfoJson) (ClassInfo, error)
}

type ClassInfoUsecase struct {
	repo   ClassInfoRepo
	logger *logger.Logger
}

func NewClassInfoUsecase(repo ClassInfoRepo, logger *logger.Logger) *ClassInfoUsecase {
	return &ClassInfoUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *ClassInfoUsecase) Create(ctx context.Context, class *ClassInfo) error {
	return uc.repo.CreateClassInfo(ctx, class)
}

func (uc *ClassInfoUsecase) DeleteByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteClassInfo(ctx, blockNumber)
}

func (uc ClassInfoUsecase) ParseClassMetadata(blockNumber uint64, classJson ClassInfoJson) (ClassInfo, error) {
	return uc.repo.ParseClassInfo(blockNumber, classJson)
}
