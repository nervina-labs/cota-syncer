package data

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	"github.com/nervina-labs/cota-smt-go/smt"
	"hash/crc32"
)

var _ biz.MintCotaKvPairRepo = (*mintCotaKvPairRepo)(nil)

type mintCotaKvPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func (rp mintCotaKvPairRepo) ParseMintCotaEntries(blockNumber uint64, entry biz.Entry) (updatedDefineCotas []biz.DefineCotaNftKvPair, withdrawCotas []biz.WithdrawCotaNftKvPair, err error) {
	entries := smt.MintCotaNFTEntriesFromSliceUnchecked(entry.Witness)
	defineCotaKeyVec := entries.DefineKeys()
	defineCotaValueVec := entries.DefineNewValues()
	lockHash, err := entry.LockScript.Hash()
	lockHashStr := lockHash.String()
	lockHashCRC32 := crc32.ChecksumIEEE([]byte(lockHashStr))
	for i := uint(0); i < defineCotaKeyVec.Len(); i++ {
		key := defineCotaKeyVec.Get(i)
		value := defineCotaValueVec.Get(i)
		if err != nil {
			return
		}
		updatedDefineCotas = append(updatedDefineCotas, biz.DefineCotaNftKvPair{
			BlockNumber: blockNumber,
			CotaId:      hex.EncodeToString(key.CotaId().RawData()),
			Total:       binary.BigEndian.Uint32(value.Total().RawData()),
			Issued:      binary.BigEndian.Uint32(value.Issued().RawData()),
			Configure:   value.Configure().AsSlice()[0],
			LockHash:    lockHashStr,
			LockHashCRC: lockHashCRC32,
		})
	}
	withdrawKeyVec := entries.WithdrawalKeys()
	withdrawValueVec := entries.WithdrawalValues()
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

func NewMintCotaKvPairRepo(data *Data, logger *logger.Logger) biz.MintCotaKvPairRepo {
	return &mintCotaKvPairRepo{
		data:   data,
		logger: logger,
	}
}
