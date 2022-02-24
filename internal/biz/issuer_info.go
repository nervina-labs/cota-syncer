package biz

import (
	"context"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
)

type IssuerInfo struct {
	ID           uint
	BlockNumber  uint64
	LockHash     string
	LockHashCRC  uint32
	Version      string
	Name         string
	Avatar       string
	Description  string
	Localization string
}

type IssuerInfoRepo interface {
	CreateIssuerInfo(ctx context.Context, issuer *IssuerInfo) error
	DeleteIssuerInfo(ctx context.Context, blockNumber uint64) error
	ParseIssuerInfo(blockNumber uint64, entry Entry) (IssuerInfo, error)
}

type IssuerInfoUsecase struct {
	repo   IssuerInfoRepo
	logger *logger.Logger
}

func NewIssuerInfoUsecase(repo IssuerInfoRepo, logger *logger.Logger) *IssuerInfoUsecase {
	return &IssuerInfoUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *IssuerInfoUsecase) Create(ctx context.Context, issuer *IssuerInfo) error {
	return uc.repo.CreateIssuerInfo(ctx, issuer)
}

func (uc *IssuerInfoUsecase) DeleteByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteIssuerInfo(ctx, blockNumber)
}

func (uc IssuerInfoUsecase) ParseIssuerMetadata(blockNumber uint64, entry Entry) (IssuerInfo, error) {
	return uc.repo.ParseIssuerInfo(blockNumber, entry)
}
