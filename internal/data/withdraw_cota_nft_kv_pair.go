package data

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	"github.com/nervina-labs/cota-smt-go/smt"
	"gorm.io/gorm"
	"hash/crc32"
)

var _ biz.WithdrawCotaNftKvPairRepo = (*withdrawCotaNftKvPairRepo)(nil)

type WithdrawCotaNftKvPair struct {
	gorm.Model

	BlockNumber         uint64
	CotaId              string
	CotaIdCRC           uint32
	TokenIndex          uint32
	OutPoint            string
	OutPointCrc         uint32
	State               uint8
	Configure           uint8
	Characteristic      string
	ReceiverLockHash    string
	ReceiverLockHashCrc uint32
	LockHash            string
	LockHashCrc         uint32
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
	lockHashStr := lockHash.String()
	lockHashCRC32 := crc32.ChecksumIEEE([]byte(lockHashStr))
	for i := uint(0); i < withdrawKeyVec.Len(); i++ {
		key := withdrawKeyVec.Get(i)
		value := withdrawValueVec.Get(i)
		cotaId := hex.EncodeToString(key.CotaId().RawData())
		outpointStr := hex.EncodeToString(value.OutPoint().RawData())
		receiverLockHashStr := hex.EncodeToString(value.To().RawData())
		withdrawCotas = append(withdrawCotas, biz.WithdrawCotaNftKvPair{
			BlockNumber:         blockNumber,
			CotaId:              cotaId,
			CotaIdCRC:           crc32.ChecksumIEEE([]byte(cotaId)),
			TokenIndex:          binary.BigEndian.Uint32(key.Index().RawData()),
			OutPoint:            outpointStr,
			OutPointCrc:         crc32.ChecksumIEEE([]byte(outpointStr)),
			State:               value.NftInfo().State().AsSlice()[0],
			Configure:           value.NftInfo().Configure().AsSlice()[0],
			Characteristic:      hex.EncodeToString(value.NftInfo().Characteristic().RawData()),
			ReceiverLockHash:    receiverLockHashStr,
			ReceiverLockHashCrc: crc32.ChecksumIEEE([]byte(receiverLockHashStr)),
			LockHash:            lockHashStr,
			LockHashCrc:         lockHashCRC32,
		})
	}
	return
}
