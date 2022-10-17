package data

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"hash/crc32"
	"time"

	"github.com/nervina-labs/cota-smt-go/smt"
	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/data/blockchain"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
)

var _ biz.RegisterCotaKvPairRepo = (*registerCotaKvPairRepo)(nil)

type RegisterCotaKvPair struct {
	ID           uint `gorm:"primaryKey"`
	BlockNumber  uint64
	LockHash     string
	CotaCellID   uint64
	LockScriptId uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
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

func (rp registerCotaKvPairRepo) ParseRegistryEntries(ctx context.Context, blockNumber uint64, tx *ckbTypes.Transaction) (registerCotas []biz.RegisterCotaKvPair, err error) {
	bytes, err := blockchain.WitnessArgsFromSliceUnchecked(tx.Witnesses[0]).InputType().IntoBytes()
	if err != nil {
		return
	}
	registerWitnessType := bytes.RawData()
	registryEntries := smt.CotaNFTRegistryEntriesFromSliceUnchecked(registerWitnessType)
	registryVec := registryEntries.Registries()
	lockMap, err := rp.generateLockMap(tx)
	if err != nil {
		return
	}
	for i := uint(0); i < registryVec.Len(); i++ {
		lockHash := hex.EncodeToString(registryVec.Get(i).LockHash().RawData())
		lock, ok := lockMap[lockHash]
		if ok {
			if err = rp.FindOrCreateScript(ctx, lock); err != nil {
				return
			}
			registerCotas = append(registerCotas, biz.RegisterCotaKvPair{
				BlockNumber:  blockNumber,
				LockHash:     lockHash,
				CotaCellID:   binary.BigEndian.Uint64(registryVec.Get(i).State().AsSlice()[0:8]),
				LockScriptId: lock.ID,
			})
		} else {
			// The outputs of the update CCID transactions have no cota cells, the lock script will be nil.
			registerCotas = append(registerCotas, biz.RegisterCotaKvPair{
				BlockNumber: blockNumber,
				LockHash:    lockHash,
				CotaCellID:  binary.BigEndian.Uint64(registryVec.Get(i).State().AsSlice()[0:8]),
			})
		}
	}
	return
}

func (rp registerCotaKvPairRepo) FindOrCreateScript(ctx context.Context, script *biz.Script) error {
	ht, err := hashType(script.HashType)
	if err != nil {
		return err
	}
	s := Script{}
	if err = rp.data.db.WithContext(ctx).FirstOrCreate(&s, Script{
		CodeHash:    script.CodeHash,
		CodeHashCrc: crc32.ChecksumIEEE([]byte(script.CodeHash)),
		HashType:    ht,
		Args:        script.Args,
		ArgsCrc:     crc32.ChecksumIEEE([]byte(script.Args)),
	}).Error; err != nil {
		return err
	}
	script.ID = s.ID
	return nil
}

func (rp registerCotaKvPairRepo) generateLockMap(tx *ckbTypes.Transaction) (map[string]*biz.Script, error) {
	lockMap := make(map[string]*biz.Script, len(tx.Outputs))
	for _, output := range tx.Outputs {
		if output.Type != nil {
			lockHash, err := output.Lock.Hash()
			if err != nil {
				return lockMap, err
			}
			hashType, err := output.Lock.HashType.Serialize()
			if err != nil {
				return lockMap, err
			}
			lock := biz.Script{
				CodeHash: hex.EncodeToString(output.Lock.CodeHash[:]),
				HashType: hex.EncodeToString(hashType),
				Args:     hex.EncodeToString(output.Lock.Args),
			}
			lockMap[lockHash.String()[2:]] = &lock
		}
	}
	return lockMap, nil
}
