package data

import (
	"context"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	"gorm.io/gorm"
)

var _ biz.DefineCotaNftKvPairRepo = (*defineCotaNftKvPairRepo)(nil)

type DefineCotaNftKvPair struct {
	gorm.Model

	BlockNumber uint64
	CotaId      string
	Total       uint32
	Issued      uint32
	Configure   uint8
	LockHash    string
	LockHashCRC uint32
}

type DefineCotaNftKvPairVersion struct {
	gorm.Model

	OldBlockNumber uint64
	BlockNumber    uint64
	CotaId         string
	Total          uint32
	OldIssued      uint32
	Issued         uint32
	Configure      uint8
	LockHash       string
	ActionType     uint8 //	0-create 1-update 2-delete
	TxIndex        uint32
}

type defineCotaNftKvPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewDefineCotaNftKvPairRepo(data *Data, logger *logger.Logger) biz.DefineCotaNftKvPairRepo {
	return &defineCotaNftKvPairRepo{
		data:   data,
		logger: logger,
	}
}

func (rp defineCotaNftKvPairRepo) CreateDefineCotaNftKvPair(ctx context.Context, d *biz.DefineCotaNftKvPair) error {
	if err := rp.data.db.WithContext(ctx).Create(d).Error; err != nil {
		return err
	}
	return nil
}

func (rp defineCotaNftKvPairRepo) DeleteDefineCotaNftKvPairs(ctx context.Context, blockNumber uint64) error {
	if err := rp.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Error; err != nil {
		return err
	}
	return nil
}
