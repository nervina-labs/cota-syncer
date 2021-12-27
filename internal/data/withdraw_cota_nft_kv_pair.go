package data

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/data/blockchain"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	"github.com/nervina-labs/cota-smt-go/smt"
	"hash/crc32"
	"time"
)

var _ biz.WithdrawCotaNftKvPairRepo = (*withdrawCotaNftKvPairRepo)(nil)

type WithdrawCotaNftKvPair struct {
	ID                   uint `gorm:"primaryKey"`
	BlockNumber          uint64
	CotaId               string
	CotaIdCRC            uint32
	TokenIndex           uint32
	OutPoint             string
	OutPointCrc          uint32
	State                uint8
	Configure            uint8
	Characteristic       string
	ReceiverLockScriptId uint
	LockHash             string
	LockHashCrc          uint32
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type Script struct {
	ID          uint `gorm:"primaryKey"`
	CodeHash    string
	CodeHashCrc uint32
	HashType    int
	Args        string
	ArgsCrc     uint32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type withdrawCotaNftKvPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewWithdrawCotaNftKvPairRepo(data *Data, logger *logger.Logger) biz.WithdrawCotaNftKvPairRepo {
	return &withdrawCotaNftKvPairRepo{
		data:   data,
		logger: logger,
	}
}

func (rp withdrawCotaNftKvPairRepo) CreateWithdrawCotaNftKvPair(ctx context.Context, w *biz.WithdrawCotaNftKvPair) error {
	if err := rp.data.db.WithContext(ctx).Create(w).Error; err != nil {
		return err
	}
	return nil
}

func (rp withdrawCotaNftKvPairRepo) DeleteWithdrawCotaNftKvPairs(ctx context.Context, blockNumber uint64) error {
	if err := rp.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(WithdrawCotaNftKvPair{}).Error; err != nil {
		return err
	}
	return nil
}

func (rp withdrawCotaNftKvPairRepo) ParseWithdrawCotaEntries(blockNumber uint64, entry biz.Entry) (withdrawCotas []biz.WithdrawCotaNftKvPair, err error) {
	entries := smt.WithdrawalCotaNFTEntriesFromSliceUnchecked(entry.Witness[1:])
	withdrawKeyVec := entries.WithdrawalKeys()
	withdrawValueVec := entries.WithdrawalValues()
	lockHash, err := entry.LockScript.Hash()
	if err != nil {
		return
	}
	lockHashStr := lockHash.String()[2:]
	lockHashCRC32 := crc32.ChecksumIEEE([]byte(lockHashStr))
	for i := uint(0); i < withdrawKeyVec.Len(); i++ {
		key := withdrawKeyVec.Get(i)
		value := withdrawValueVec.Get(i)
		cotaId := hex.EncodeToString(key.CotaId().RawData())
		outpointStr := hex.EncodeToString(value.OutPoint().RawData())
		receiverLock := blockchain.ScriptFromSliceUnchecked(value.ToLock().RawData())
		script := biz.Script{
			CodeHash: hex.EncodeToString(receiverLock.CodeHash().RawData()),
			HashType: hex.EncodeToString(receiverLock.HashType().AsSlice()),
			Args:     hex.EncodeToString(receiverLock.Args().RawData()),
		}
		err = rp.FindOrCreateScript(context.TODO(), script)
		if err != nil {
			return nil, err
		}
		withdrawCotas = append(withdrawCotas, biz.WithdrawCotaNftKvPair{
			BlockNumber:          blockNumber,
			CotaId:               cotaId,
			CotaIdCRC:            crc32.ChecksumIEEE([]byte(cotaId)),
			TokenIndex:           binary.BigEndian.Uint32(key.Index().RawData()),
			OutPoint:             outpointStr,
			OutPointCrc:          crc32.ChecksumIEEE([]byte(outpointStr)),
			State:                value.NftInfo().State().AsSlice()[0],
			Configure:            value.NftInfo().Configure().AsSlice()[0],
			Characteristic:       hex.EncodeToString(value.NftInfo().Characteristic().RawData()),
			ReceiverLockScriptId: script.ID,
			LockHash:             lockHashStr,
			LockHashCrc:          lockHashCRC32,
		})
	}
	return
}

func (rp withdrawCotaNftKvPairRepo) FindOrCreateScript(ctx context.Context, script biz.Script) error {
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

func hashType(hashTypeStr string) (int, error) {
	switch hashTypeStr {
	case "data":
		return 0, nil
	case "type":
		return 1, nil
	case "data1":
		return 2, nil
	}
	return -1, errors.New("invalid hash type")
}
