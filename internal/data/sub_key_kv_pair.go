package data

import (
	"context"
	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	"time"
)

type SubKeyKvPair struct {
	ID          uint `gorm:"primaryKey"`
	BlockNumber uint64
	LockHash    string
	SubType     string
	ExtData     uint32
	AlgIndex    uint16
	PubkeyHash  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SubKeyKvPairVersion struct {
	ID             uint `gorm:"primaryKey"`
	OldBlockNumber uint64
	BlockNumber    uint64
	LockHash       string
	SubType        string
	ExtData        uint32
	OldAlgIndex    uint16
	AlgIndex       uint16
	OldPubkeyHash  string
	PubkeyHash     string
	ActionType     uint8 //	0-create 1-update
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

var _ biz.SubKeyPairRepo = (*subKeyPairRepo)(nil)

type subKeyPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewSubKeyKvPairRepo(data *Data, logger *logger.Logger) biz.SubKeyPairRepo {
	return &subKeyPairRepo{
		data:   data,
		logger: logger,
	}
}

func (rp subKeyPairRepo) CreateSubKeyPair(ctx context.Context, subKey *biz.SubKeyPair) error {
	if err := rp.data.db.WithContext(ctx).Create(subKey).Error; err != nil {
		return err
	}
	return nil
}

func (rp subKeyPairRepo) DeleteSubKeyPairs(ctx context.Context, blockNumber uint64) error {
	if err := rp.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(SubKeyKvPair{}).Error; err != nil {
		return err
	}
	return nil
}
