package data

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"hash/crc32"
	"strconv"
	"time"

	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/data/blockchain"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	"github.com/nervina-labs/cota-smt-go/smt"
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
	TxHash               string
	State                uint8
	Configure            uint8
	Characteristic       string
	ReceiverLockScriptId uint
	LockHash             string
	LockHashCrc          uint32
	LockScriptId         uint
	Version              uint8
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type Script struct {
	ID          uint `gorm:"primaryKey"`
	CodeHash    string
	CodeHashCrc uint32
	HashType    int64
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

func (rp withdrawCotaNftKvPairRepo) ParseWithdrawCotaEntries(blockNumber uint64, entry biz.Entry) ([]biz.WithdrawCotaNftKvPair, error) {
	if entry.Version == 0 {
		return generateV0WithdrawKvPair(blockNumber, entry, rp)
	}
	return generateV1WithdrawKvPair(blockNumber, entry, rp)
}

func (rp withdrawCotaNftKvPairRepo) FindOrCreateScript(ctx context.Context, script *biz.Script) error {
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

func hashType(hashTypeStr string) (int64, error) {
	return strconv.ParseInt(hashTypeStr, 16, 32)
}

func generateV0WithdrawKvPair(blockNumber uint64, entry biz.Entry, rp withdrawCotaNftKvPairRepo) ([]biz.WithdrawCotaNftKvPair, error) {
	var withdrawCotas []biz.WithdrawCotaNftKvPair
	entries, err := smt.WithdrawalCotaNFTEntriesFromSlice(entry.InputType[1:], false)
	if err != nil {
		return withdrawCotas, err
	}
	withdrawKeyVec := entries.WithdrawalKeys()
	withdrawValueVec := entries.WithdrawalValues()
	senderLock, lockHashStr, lockHashCRC32, err := generateSenderLock(entry, rp)
	if err != nil {
		return withdrawCotas, err
	}
	for i := uint(0); i < withdrawKeyVec.Len(); i++ {
		key := withdrawKeyVec.Get(i)
		value := withdrawValueVec.Get(i)
		cotaId := hex.EncodeToString(key.CotaId().RawData())
		outpointStr := hex.EncodeToString(value.OutPoint().RawData())
		receiverLock, err := generateReceiverLock(value.ToLock().RawData(), rp)
		if err != nil {
			return withdrawCotas, err
		}
		withdrawCotas = append(withdrawCotas, biz.WithdrawCotaNftKvPair{
			BlockNumber:          blockNumber,
			CotaId:               cotaId,
			CotaIdCRC:            crc32.ChecksumIEEE([]byte(cotaId)),
			TokenIndex:           binary.BigEndian.Uint32(key.Index().RawData()),
			OutPoint:             outpointStr,
			OutPointCrc:          crc32.ChecksumIEEE([]byte(outpointStr)),
			TxHash:               entry.TxHash.String()[2:],
			State:                value.NftInfo().State().AsSlice()[0],
			Configure:            value.NftInfo().Configure().AsSlice()[0],
			Characteristic:       hex.EncodeToString(value.NftInfo().Characteristic().RawData()),
			ReceiverLockScriptId: receiverLock.ID,
			LockHash:             lockHashStr,
			LockHashCrc:          lockHashCRC32,
			LockScriptId:         senderLock.ID,
			Version:              entry.Version,
		})
	}
	return withdrawCotas, nil
}

func generateV1WithdrawKvPair(blockNumber uint64, entry biz.Entry, rp withdrawCotaNftKvPairRepo) ([]biz.WithdrawCotaNftKvPair, error) {
	var withdrawCotas []biz.WithdrawCotaNftKvPair
	entries, err := smt.WithdrawalCotaNFTV1EntriesFromSlice(entry.InputType[1:], false)
	if err != nil {
		return withdrawCotas, err
	}
	withdrawKeyVec := entries.WithdrawalKeys()
	withdrawValueVec := entries.WithdrawalValues()
	senderLock, lockHashStr, lockHashCRC32, err := generateSenderLock(entry, rp)
	if err != nil {
		return withdrawCotas, err
	}
	for i := uint(0); i < withdrawKeyVec.Len(); i++ {
		key := withdrawKeyVec.Get(i)
		value := withdrawValueVec.Get(i)
		cotaId := hex.EncodeToString(key.NftId().CotaId().RawData())
		outpointStr := hex.EncodeToString(key.OutPoint().RawData())
		receiverLock, err := generateReceiverLock(value.ToLock().RawData(), rp)
		if err != nil {
			return withdrawCotas, err
		}
		withdrawCotas = append(withdrawCotas, biz.WithdrawCotaNftKvPair{
			BlockNumber:          blockNumber,
			CotaId:               cotaId,
			CotaIdCRC:            crc32.ChecksumIEEE([]byte(cotaId)),
			TokenIndex:           binary.BigEndian.Uint32(key.NftId().Index().RawData()),
			OutPoint:             outpointStr,
			OutPointCrc:          crc32.ChecksumIEEE([]byte(outpointStr)),
			TxHash:               entry.TxHash.String()[2:],
			State:                value.NftInfo().State().AsSlice()[0],
			Configure:            value.NftInfo().Configure().AsSlice()[0],
			Characteristic:       hex.EncodeToString(value.NftInfo().Characteristic().RawData()),
			ReceiverLockScriptId: receiverLock.ID,
			LockHash:             lockHashStr,
			LockHashCrc:          lockHashCRC32,
			LockScriptId:         senderLock.ID,
			Version:              entry.Version,
		})
	}
	return withdrawCotas, nil
}

func generateSenderLock(entry biz.Entry, rp withdrawCotaNftKvPairRepo) (biz.Script, string, uint32, error) {
	var (
		lockScript    biz.Script
		lockHashStr   string
		lockHashCRC32 uint32
	)
	lockScript = biz.Script{
		CodeHash: entry.LockScript.CodeHash.String()[2:],
		HashType: string(entry.LockScript.HashType),
		Args:     hex.EncodeToString(entry.LockScript.Args),
	}
	if err := rp.FindOrCreateScript(context.TODO(), &lockScript); err != nil {
		return lockScript, lockHashStr, lockHashCRC32, err
	}
	lockHash, err := entry.LockScript.Hash()
	if err != nil {
		return lockScript, lockHashStr, lockHashCRC32, err
	}
	lockHashStr = lockHash.String()[2:]
	lockHashCRC32 = crc32.ChecksumIEEE([]byte(lockHashStr))
	return lockScript, lockHashStr, lockHashCRC32, nil
}

func generateReceiverLock(slice []byte, rp withdrawCotaNftKvPairRepo) (biz.Script, error) {
	receiverLock := blockchain.ScriptFromSliceUnchecked(slice)
	script := biz.Script{
		CodeHash: hex.EncodeToString(receiverLock.CodeHash().RawData()),
		HashType: hex.EncodeToString(receiverLock.HashType().AsSlice()),
		Args:     hex.EncodeToString(receiverLock.Args().RawData()),
	}
	err := rp.FindOrCreateScript(context.TODO(), &script)
	return script, err
}
