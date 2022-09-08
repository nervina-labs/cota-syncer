package biz

import (
	"context"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
)

type IssuerInfo struct {
	BlockNumber  uint64
	LockHash     string
	Version      string
	Name         string
	Avatar       string
	Description  string
	Localization string
	TxIndex      uint32
}

type IssuerInfoRepo interface {
	CreateIssuerInfo(ctx context.Context, issuer *IssuerInfo) error
	DeleteIssuerInfo(ctx context.Context, blockNumber uint64) error
	ParseIssuerInfo(blockNumber uint64, txIndex uint32, lockScript *ckbTypes.Script, issuerMeta map[string]any) (IssuerInfo, error)
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

func (uc IssuerInfoUsecase) ParseMetadata(blockNumber uint64, txIndex uint32, lockScript *ckbTypes.Script, issuerMeta map[string]any) (IssuerInfo, error) {
	return uc.repo.ParseIssuerInfo(blockNumber, txIndex, lockScript, issuerMeta)
}
