package data

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/data/blockchain"
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
	if entry.Version == 0 {
		err = generateV0KvPairs(blockNumber, entry, updatedDefineCotas, rp, withdrawCotas)
		return
	}
	err = generateV1KvPairs(blockNumber, entry, updatedDefineCotas, rp, withdrawCotas)
	return
}

func generateV1KvPairs(blockNumber uint64, entry biz.Entry, updatedDefineCotas []biz.DefineCotaNftKvPair, rp mintCotaKvPairRepo, withdrawCotas []biz.WithdrawCotaNftKvPair) error {
	entries := smt.MintCotaNFTV1EntriesFromSliceUnchecked(entry.Witness[1:])
	defineCotaKeyVec := entries.DefineKeys()
	defineCotaValueVec := entries.DefineNewValues()
	lockHash, err := entry.LockScript.Hash()
	if err != nil {
		return err
	}
	lockHashStr := lockHash.String()[2:]
	lockHashCRC32 := crc32.ChecksumIEEE([]byte(lockHashStr))
	for i := uint(0); i < defineCotaKeyVec.Len(); i++ {
		key := defineCotaKeyVec.Get(i)
		value := defineCotaValueVec.Get(i)
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
		cotaId := hex.EncodeToString(key.NftId().CotaId().RawData())
		outpointStr := hex.EncodeToString(key.OutPoint().RawData())
		receiverLock := blockchain.ScriptFromSliceUnchecked(value.ToLock().RawData())
		script := biz.Script{
			CodeHash: hex.EncodeToString(receiverLock.CodeHash().RawData()),
			HashType: hex.EncodeToString(receiverLock.HashType().AsSlice()),
			Args:     hex.EncodeToString(receiverLock.Args().RawData()),
		}
		err = rp.FindOrCreateScript(context.TODO(), &script)
		if err != nil {
			return err
		}
		withdrawCotas = append(withdrawCotas, biz.WithdrawCotaNftKvPair{
			BlockNumber:          blockNumber,
			CotaId:               cotaId,
			CotaIdCRC:            crc32.ChecksumIEEE([]byte(cotaId)),
			TokenIndex:           binary.BigEndian.Uint32(key.NftId().Index().RawData()),
			OutPoint:             outpointStr,
			OutPointCrc:          crc32.ChecksumIEEE([]byte(outpointStr)),
			State:                value.NftInfo().State().AsSlice()[0],
			Configure:            value.NftInfo().Configure().AsSlice()[0],
			Characteristic:       hex.EncodeToString(value.NftInfo().Characteristic().RawData()),
			ReceiverLockScriptId: script.ID,
			LockHash:             lockHashStr,
			LockHashCrc:          lockHashCRC32,
			Version:              entry.Version,
		})
	}
	return nil
}

func generateV0KvPairs(blockNumber uint64, entry biz.Entry, updatedDefineCotas []biz.DefineCotaNftKvPair, rp mintCotaKvPairRepo, withdrawCotas []biz.WithdrawCotaNftKvPair) error {
	entries := smt.MintCotaNFTEntriesFromSliceUnchecked(entry.Witness[1:])
	defineCotaKeyVec := entries.DefineKeys()
	defineCotaValueVec := entries.DefineNewValues()
	lockHash, err := entry.LockScript.Hash()
	if err != nil {
		return err
	}
	lockHashStr := lockHash.String()[2:]
	lockHashCRC32 := crc32.ChecksumIEEE([]byte(lockHashStr))
	for i := uint(0); i < defineCotaKeyVec.Len(); i++ {
		key := defineCotaKeyVec.Get(i)
		value := defineCotaValueVec.Get(i)
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
		receiverLock := blockchain.ScriptFromSliceUnchecked(value.ToLock().RawData())
		script := biz.Script{
			CodeHash: hex.EncodeToString(receiverLock.CodeHash().RawData()),
			HashType: hex.EncodeToString(receiverLock.HashType().AsSlice()),
			Args:     hex.EncodeToString(receiverLock.Args().RawData()),
		}
		err = rp.FindOrCreateScript(context.TODO(), &script)
		if err != nil {
			return err
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
			Version:              entry.Version,
		})
	}
	return nil
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
