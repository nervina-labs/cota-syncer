package data

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"hash/crc32"
	"time"

	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	"github.com/nervina-labs/cota-smt-go/smt"
)

var _ biz.MintCotaKvPairRepo = (*mintCotaKvPairRepo)(nil)

type mintCotaKvPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func (rp mintCotaKvPairRepo) ParseMintCotaEntries(blockNumber uint64, entry biz.Entry) ([]biz.DefineCotaNftKvPair, []biz.WithdrawCotaNftKvPair, error) {
	if entry.Version == 0 {
		return generateMintV0KvPairs(blockNumber, entry, rp)
	}
	return generateMintV1KvPairs(blockNumber, entry, rp)
}

func generateMintV1KvPairs(blockNumber uint64, entry biz.Entry, rp mintCotaKvPairRepo) ([]biz.DefineCotaNftKvPair, []biz.WithdrawCotaNftKvPair, error) {
	var (
		defineCotas   []biz.DefineCotaNftKvPair
		withdrawCotas []biz.WithdrawCotaNftKvPair
	)
	entries, err := smt.MintCotaNFTV1EntriesFromSlice(entry.InputType[1:], false)
	if err != nil {
		return defineCotas, withdrawCotas, err
	}

	defineCotaKeyVec := entries.DefineKeys()
	defineCotaValueVec := entries.DefineNewValues()
	senderLock, lockHashStr, lockHashCRC32, err := GenerateSenderLock(entry)
	if err != nil {
		return defineCotas, withdrawCotas, err
	}
	if err := rp.FindOrCreateScript(context.TODO(), &senderLock); err != nil {
		return defineCotas, withdrawCotas, err
	}
	for i := uint(0); i < defineCotaKeyVec.Len(); i++ {
		key := defineCotaKeyVec.Get(i)
		value := defineCotaValueVec.Get(i)
		defineCotas = append(defineCotas, biz.DefineCotaNftKvPair{
			BlockNumber: blockNumber,
			CotaId:      hex.EncodeToString(key.CotaId().RawData()),
			Total:       binary.BigEndian.Uint32(value.Total().RawData()),
			Issued:      binary.BigEndian.Uint32(value.Issued().RawData()),
			Configure:   value.Configure().AsSlice()[0],
			LockHash:    lockHashStr,
			LockHashCRC: lockHashCRC32,
			UpdatedAt:   time.Now(),
		})
	}
	withdrawKeyVec := entries.WithdrawalKeys()
	withdrawValueVec := entries.WithdrawalValues()
	for i := uint(0); i < withdrawKeyVec.Len(); i++ {
		key := withdrawKeyVec.Get(i)
		value := withdrawValueVec.Get(i)
		cotaId := hex.EncodeToString(key.NftId().CotaId().RawData())
		outpointStr := hex.EncodeToString(key.OutPoint().RawData())
		receiverLock := GenerateReceiverLock(value.ToLock().RawData())
		err = rp.FindOrCreateScript(context.TODO(), &receiverLock)
		if err != nil {
			return defineCotas, withdrawCotas, err
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
	return defineCotas, withdrawCotas, err
}

func generateMintV0KvPairs(blockNumber uint64, entry biz.Entry, rp mintCotaKvPairRepo) ([]biz.DefineCotaNftKvPair, []biz.WithdrawCotaNftKvPair, error) {
	var (
		defineCotas   []biz.DefineCotaNftKvPair
		withdrawCotas []biz.WithdrawCotaNftKvPair
	)
	entries, err := smt.MintCotaNFTEntriesFromSlice(entry.InputType[1:], false)
	if err != nil {
		return defineCotas, withdrawCotas, err
	}
	defineCotaKeyVec := entries.DefineKeys()
	defineCotaValueVec := entries.DefineNewValues()
	senderLock, lockHashStr, lockHashCRC32, err := GenerateSenderLock(entry)
	if err != nil {
		return defineCotas, withdrawCotas, err
	}
	if err := rp.FindOrCreateScript(context.TODO(), &senderLock); err != nil {
		return defineCotas, withdrawCotas, err
	}
	for i := uint(0); i < defineCotaKeyVec.Len(); i++ {
		key := defineCotaKeyVec.Get(i)
		value := defineCotaValueVec.Get(i)
		defineCotas = append(defineCotas, biz.DefineCotaNftKvPair{
			BlockNumber: blockNumber,
			CotaId:      hex.EncodeToString(key.CotaId().RawData()),
			Total:       binary.BigEndian.Uint32(value.Total().RawData()),
			Issued:      binary.BigEndian.Uint32(value.Issued().RawData()),
			Configure:   value.Configure().AsSlice()[0],
			LockHash:    lockHashStr,
			LockHashCRC: lockHashCRC32,
			UpdatedAt:   time.Now().UTC(),
		})
	}
	withdrawKeyVec := entries.WithdrawalKeys()
	withdrawValueVec := entries.WithdrawalValues()
	for i := uint(0); i < withdrawKeyVec.Len(); i++ {
		key := withdrawKeyVec.Get(i)
		value := withdrawValueVec.Get(i)
		cotaId := hex.EncodeToString(key.CotaId().RawData())
		outpointStr := hex.EncodeToString(value.OutPoint().RawData())
		receiverLock := GenerateReceiverLock(value.ToLock().RawData())
		err = rp.FindOrCreateScript(context.TODO(), &receiverLock)
		if err != nil {
			return defineCotas, withdrawCotas, err
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
	return defineCotas, withdrawCotas, nil
}

func (rp mintCotaKvPairRepo) FindOrCreateScript(ctx context.Context, script *biz.Script) error {
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

func NewMintCotaKvPairRepo(data *Data, logger *logger.Logger) biz.MintCotaKvPairRepo {
	return &mintCotaKvPairRepo{
		data:   data,
		logger: logger,
	}
}
