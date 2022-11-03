package data

import (
	"context"
	"github.com/nervina-labs/cota-smt-go/smt"
	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strconv"
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
	if err := rp.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Error; err != nil {
		return err
	}
	return nil
}

func (rp subKeyPairRepo) ParseSubKeyPairs(blockNumber uint64, entry biz.Entry) ([]biz.SubKeyPair, error) {
	var (
		pairs             []biz.SubKeyPair
		lockHash          types.Hash
		extData, algIndex int64
		err               error
	)

	extensionEntries := smt.ExtensionEntriesFromSliceUnchecked(entry.InputType[1:])
	if string(extensionEntries.SubType().RawData()) != "subkey" {
		return []biz.SubKeyPair{}, nil
	}

	entries := smt.SubKeyEntriesFromSliceUnchecked(extensionEntries.RawData().AsSlice())
	if lockHash, err = entry.LockScript.Hash(); err != nil {
		return nil, err
	}

	for i := uint(0); i < entries.Len(); i++ {
		key := entries.Keys().Get(i)
		value := entries.Values().Get(i)

		if extData, err = strconv.ParseInt(string(key.ExtData().RawData()), 10, 64); err != nil {
			return nil, err
		}
		if algIndex, err = strconv.ParseInt(string(value.AlgIndex().RawData()), 10, 64); err != nil {
			return nil, err
		}

		pairs = append(pairs, biz.SubKeyPair{
			BlockNumber: blockNumber,
			LockHash:    remove0x(lockHash.Hex()),
			SubType:     string(key.SubType().RawData()),
			ExtData:     uint32(extData),
			AlgIndex:    uint16(algIndex),
			PubkeyHash:  remove0x(string(value.PubkeyHash().RawData())),
			UpdatedAt:   time.Now().UTC(),
		})
	}

	return pairs, nil
}
