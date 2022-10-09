package biz

import (
	"context"

	"github.com/nervina-labs/cota-syncer/internal/logger"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
)

type SubKeyInfo struct {
	BlockNumber  uint64
	LockHash     string
	PubKey       string
	CredentialId string
	Alg          string
}

type JoyIDInfo struct {
	BlockNumber  uint64
	LockHash     string
	Version      string
	PubKey       string
	CredentialId string
	Alg          string
	CotaCellId   string
	Name         string
	Avatar       string
	Description  string
	Extension    string
	Nickname     string
	TxIndex      uint32
	SubKeys      []SubKeyInfo
}

type JoyIDInfoRepo interface {
	CreateJoyIDInfo(ctx context.Context, joyID *JoyIDInfo) error
	DeleteJoyIDInfo(ctx context.Context, blockNumber uint64) error
	ParseJoyIDInfo(ctx context.Context, blockNumber uint64, txIndex uint32, lockScript *ckbTypes.Script, joyIDMeta map[string]any) (JoyIDInfo, error)
}

type JoyIDInfoUsecase struct {
	repo   JoyIDInfoRepo
	logger *logger.Logger
}

func NewJoyIDInfoUsecase(repo JoyIDInfoRepo, logger *logger.Logger) *JoyIDInfoUsecase {
	return &JoyIDInfoUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *JoyIDInfoUsecase) Create(ctx context.Context, joyID *JoyIDInfo) error {
	return uc.repo.CreateJoyIDInfo(ctx, joyID)
}

func (uc *JoyIDInfoUsecase) DeleteByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return uc.repo.DeleteJoyIDInfo(ctx, blockNumber)
}

func (uc JoyIDInfoUsecase) ParseMetadata(ctx context.Context, blockNumber uint64, txIndex uint32, lockScript *ckbTypes.Script, joyIDMeta map[string]any) (JoyIDInfo, error) {
	return uc.repo.ParseJoyIDInfo(ctx, blockNumber, txIndex, lockScript, joyIDMeta)
}
