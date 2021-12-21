package data

import (
	"context"
	"encoding/hex"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/data/blockchain"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	"github.com/nervina-labs/cota-smt-go/smt"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
	"time"
)

var _ biz.RegisterCotaKvPairRepo = (*registerCotaKvPairRepo)(nil)

type RegisterCotaKvPair struct {
	ID          uint `gorm:"primaryKey"`
	BlockNumber uint64
	LockHash    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type registerCotaKvPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewRegisterCotaKvPairRepo(data *Data, logger *logger.Logger) biz.RegisterCotaKvPairRepo {
	return &registerCotaKvPairRepo{
		data:   data,
		logger: logger,
	}
}

func (rp registerCotaKvPairRepo) CreateRegisterCotaKvPair(ctx context.Context, r *biz.RegisterCotaKvPair) error {
	if err := rp.data.db.WithContext(ctx).Create(r).Error; err != nil {
		return err
	}
	return nil
}

func (rp registerCotaKvPairRepo) DeleteRegisterCotaKvPairs(ctx context.Context, blockNumber uint64) error {
	if err := rp.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(RegisterCotaKvPair{}).Error; err != nil {
		return err
	}
	return nil
}

func (rp registerCotaKvPairRepo) ParseRegistryEntries(_ context.Context, blockNumber uint64, tx *ckbTypes.Transaction) ([]biz.RegisterCotaKvPair, error) {
	bytes, err := blockchain.WitnessArgsFromSliceUnchecked(tx.Witnesses[0]).InputType().IntoBytes()
	if err != nil {
		return []biz.RegisterCotaKvPair{}, err
	}
	registerWitnessType := bytes.RawData()
	registryEntries := smt.CotaNFTRegistryEntriesFromSliceUnchecked(registerWitnessType)
	rp.logger.Infof(context.TODO(), "entries: %v", registryEntries)
	registryVec := registryEntries.Registries()
	registerCotas := make([]biz.RegisterCotaKvPair, registryVec.Len())
	for i := uint(0); i < registryVec.Len(); i++ {
		registryEntry := registryVec.Get(i)
		registerCotas[i] = biz.RegisterCotaKvPair{
			BlockNumber: blockNumber,
			LockHash:    hex.EncodeToString(registryEntry.LockHash().RawData()),
		}
	}
	return registerCotas, nil
}
