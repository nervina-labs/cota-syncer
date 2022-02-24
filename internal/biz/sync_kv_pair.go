package biz

import (
	"context"

	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
)

type KvPair struct {
	Registers          []RegisterCotaKvPair
	DefineCotas        []DefineCotaNftKvPair
	UpdatedDefineCotas []DefineCotaNftKvPair
	HoldCotas          []HoldCotaNftKvPair
	UpdatedHoldCotas   []HoldCotaNftKvPair
	WithdrawCotas      []WithdrawCotaNftKvPair
	ClaimedCotas       []ClaimedCotaNftKvPair
	IssuerInfos        []IssuerInfo
}

func (p KvPair) HasRegisters() bool {
	if len(p.Registers) > 0 {
		return true
	}
	return false
}

func (p KvPair) HasDefineCotas() bool {
	if len(p.DefineCotas) > 0 {
		return true
	}
	return false
}

func (p KvPair) HasUpdatedDefineCotas() bool {
	if len(p.UpdatedDefineCotas) > 0 {
		return true
	}
	return false
}
func (p KvPair) HasHoldCotas() bool {
	if len(p.HoldCotas) > 0 {
		return true
	}
	return false
}
func (p KvPair) HasUpdatedHoldCotas() bool {
	if len(p.UpdatedHoldCotas) > 0 {
		return true
	}
	return false
}
func (p KvPair) HasWithdrawCotas() bool {
	if len(p.WithdrawCotas) > 0 {
		return true
	}
	return false
}
func (p KvPair) HasClaimedCotas() bool {
	if len(p.ClaimedCotas) > 0 {
		return true
	}
	return false
}

func (p KvPair) HasIssuerInfos() bool {
	if len(p.IssuerInfos) > 0 {
		return true
	}
	return false
}

type KvPairRepo interface {
	CreateKvPairs(ctx context.Context, checkInfo CheckInfo, kvPair *KvPair) error
	RestoreKvPairs(ctx context.Context, blockNumber uint64) error
}

type SyncKvPairUsecase struct {
	repo   KvPairRepo
	logger *logger.Logger
}

func NewSyncKvPairUsecase(repo KvPairRepo, logger *logger.Logger) *SyncKvPairUsecase {
	return &SyncKvPairUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc SyncKvPairUsecase) CreateKvPairs(ctx context.Context, checkInfo CheckInfo, kvPair *KvPair) error {
	return uc.repo.CreateKvPairs(ctx, checkInfo, kvPair)
}

func (uc SyncKvPairUsecase) RestoreKvPairs(ctx context.Context, blockNumber uint64) error {
	return uc.repo.RestoreKvPairs(ctx, blockNumber)
}
