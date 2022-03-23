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
	ClassInfos         []ClassInfo
}

func (p KvPair) HasRegisters() bool {
	return len(p.Registers) > 0
}

func (p KvPair) HasDefineCotas() bool {
	return len(p.DefineCotas) > 0
}

func (p KvPair) HasUpdatedDefineCotas() bool {
	return len(p.UpdatedDefineCotas) > 0
}
func (p KvPair) HasHoldCotas() bool {
	return len(p.HoldCotas) > 0
}
func (p KvPair) HasUpdatedHoldCotas() bool {
	return len(p.UpdatedHoldCotas) > 0
}
func (p KvPair) HasWithdrawCotas() bool {
	return len(p.WithdrawCotas) > 0
}
func (p KvPair) HasClaimedCotas() bool {
	return len(p.ClaimedCotas) > 0
}

func (p KvPair) HasIssuerInfos() bool {
	return len(p.IssuerInfos) > 0
}

func (p KvPair) HasClassInfos() bool {
	return len(p.ClassInfos) > 0
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
