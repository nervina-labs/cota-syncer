package data

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"hash/crc32"

	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	"github.com/nervina-labs/cota-smt-go/smt"
)

var _ biz.TransferCotaKvPairRepo = (*transferCotaKvPairRepo)(nil)

type transferCotaKvPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func (rp transferCotaKvPairRepo) ParseTransferCotaEntries(blockNumber uint64, entry biz.Entry) ([]biz.ClaimedCotaNftKvPair, []biz.WithdrawCotaNftKvPair, error) {
	if entry.Version == 0 {
		return generateTransferV0KvPairs(blockNumber, entry, rp)
	}
	return generateTransferV1ToV2KvPairs(blockNumber, entry, rp)
}

func (rp transferCotaKvPairRepo) ParseTransferUpdateCotaEntries(blockNumber uint64, entry biz.Entry) ([]biz.ClaimedCotaNftKvPair, []biz.WithdrawCotaNftKvPair, error) {
	if entry.Version == 0 {
		return generateTransferUpdateV0KvPairs(blockNumber, entry, rp)
	}
	return generateTransferUpdateV1ToV2KvPairs(blockNumber, entry, rp)
}

func (rp transferCotaKvPairRepo) FindOrCreateScript(ctx context.Context, script *biz.Script) error {
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

func NewTransferCotaKvPairRepo(data *Data, logger *logger.Logger) biz.TransferCotaKvPairRepo {
	return &transferCotaKvPairRepo{
		data:   data,
		logger: logger,
	}
}

func generateTransferUpdateV0KvPairs(blockNumber uint64, entry biz.Entry, rp transferCotaKvPairRepo) ([]biz.ClaimedCotaNftKvPair, []biz.WithdrawCotaNftKvPair, error) {
	var (
		claimedCotas  []biz.ClaimedCotaNftKvPair
		withdrawCotas []biz.WithdrawCotaNftKvPair
	)
	entries, err := smt.TransferUpdateCotaNFTEntriesFromSlice(entry.InputType[1:], false)
	if err != nil {
		return claimedCotas, withdrawCotas, err
	}
	claimedCotaKeyVec := entries.ClaimKeys()
	senderLock, lockHashStr, lockHashCRC32, err := GenerateSenderLock(entry)
	if err != nil {
		return claimedCotas, withdrawCotas, err
	}
	if err := rp.FindOrCreateScript(context.TODO(), &senderLock); err != nil {
		return claimedCotas, withdrawCotas, err
	}
	for i := uint(0); i < claimedCotaKeyVec.Len(); i++ {
		key := claimedCotaKeyVec.Get(i)
		cotaId := hex.EncodeToString(key.NftId().CotaId().RawData())
		outpointStr := hex.EncodeToString(key.OutPoint().RawData())
		claimedCotas = append(claimedCotas, biz.ClaimedCotaNftKvPair{
			BlockNumber: blockNumber,
			CotaId:      hex.EncodeToString(key.NftId().CotaId().RawData()),
			CotaIdCRC:   crc32.ChecksumIEEE([]byte(cotaId)),
			TokenIndex:  binary.BigEndian.Uint32(key.NftId().Index().RawData()),
			OutPoint:    outpointStr,
			OutPointCrc: crc32.ChecksumIEEE([]byte(outpointStr)),
			LockHash:    lockHashStr,
			LockHashCrc: lockHashCRC32,
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
			return claimedCotas, withdrawCotas, err
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
	return claimedCotas, withdrawCotas, nil
}

func generateTransferUpdateV1ToV2KvPairs(blockNumber uint64, entry biz.Entry, rp transferCotaKvPairRepo) ([]biz.ClaimedCotaNftKvPair, []biz.WithdrawCotaNftKvPair, error) {
	var (
		claimedCotaKeyVec *smt.ClaimCotaNFTKeyVec
		withdrawKeyVec    *smt.WithdrawalCotaNFTKeyV1Vec
		withdrawValueVec  *smt.WithdrawalCotaNFTValueV1Vec
		claimedCotas      []biz.ClaimedCotaNftKvPair
		withdrawCotas     []biz.WithdrawCotaNftKvPair
	)

	if entry.Version == 1 {
		entries, err := smt.TransferUpdateCotaNFTV1EntriesFromSlice(entry.InputType[1:], false)
		if err != nil {
			return claimedCotas, withdrawCotas, err
		}
		claimedCotaKeyVec = entries.ClaimKeys()
		withdrawKeyVec = entries.WithdrawalKeys()
		withdrawValueVec = entries.WithdrawalValues()
	} else {
		entries, err := smt.TransferUpdateCotaNFTV2EntriesFromSlice(entry.InputType[1:], false)
		if err != nil {
			return claimedCotas, withdrawCotas, err
		}
		claimedCotaKeyVec = entries.ClaimKeys()
		withdrawKeyVec = entries.WithdrawalKeys()
		withdrawValueVec = entries.WithdrawalValues()
	}
	senderLock, lockHashStr, lockHashCRC32, err := GenerateSenderLock(entry)
	if err != nil {
		return claimedCotas, withdrawCotas, err
	}
	if err := rp.FindOrCreateScript(context.TODO(), &senderLock); err != nil {
		return claimedCotas, withdrawCotas, err
	}
	for i := uint(0); i < claimedCotaKeyVec.Len(); i++ {
		key := claimedCotaKeyVec.Get(i)
		cotaId := hex.EncodeToString(key.NftId().CotaId().RawData())
		outpointStr := hex.EncodeToString(key.OutPoint().RawData())
		claimedCotas = append(claimedCotas, biz.ClaimedCotaNftKvPair{
			BlockNumber: blockNumber,
			CotaId:      hex.EncodeToString(key.NftId().CotaId().RawData()),
			CotaIdCRC:   crc32.ChecksumIEEE([]byte(cotaId)),
			TokenIndex:  binary.BigEndian.Uint32(key.NftId().Index().RawData()),
			OutPoint:    outpointStr,
			OutPointCrc: crc32.ChecksumIEEE([]byte(outpointStr)),
			LockHash:    lockHashStr,
			LockHashCrc: lockHashCRC32,
		})
	}
	for i := uint(0); i < withdrawKeyVec.Len(); i++ {
		key := withdrawKeyVec.Get(i)
		value := withdrawValueVec.Get(i)
		cotaId := hex.EncodeToString(key.NftId().CotaId().RawData())
		outpointStr := hex.EncodeToString(key.OutPoint().RawData())
		receiverLock := GenerateReceiverLock(value.ToLock().RawData())
		err = rp.FindOrCreateScript(context.TODO(), &receiverLock)
		if err != nil {
			return claimedCotas, withdrawCotas, err
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
	return claimedCotas, withdrawCotas, nil
}

func generateTransferV0KvPairs(blockNumber uint64, entry biz.Entry, rp transferCotaKvPairRepo) ([]biz.ClaimedCotaNftKvPair, []biz.WithdrawCotaNftKvPair, error) {
	var (
		claimedCotas  []biz.ClaimedCotaNftKvPair
		withdrawCotas []biz.WithdrawCotaNftKvPair
	)
	entries, err := smt.TransferCotaNFTEntriesFromSlice(entry.InputType[1:], false)
	if err != nil {
		return claimedCotas, withdrawCotas, err
	}
	claimedCotaKeyVec := entries.ClaimKeys()
	senderLock, lockHashStr, lockHashCRC32, err := GenerateSenderLock(entry)
	if err != nil {
		return claimedCotas, withdrawCotas, err
	}
	if err := rp.FindOrCreateScript(context.TODO(), &senderLock); err != nil {
		return claimedCotas, withdrawCotas, err
	}
	for i := uint(0); i < claimedCotaKeyVec.Len(); i++ {
		key := claimedCotaKeyVec.Get(i)
		cotaId := hex.EncodeToString(key.NftId().CotaId().RawData())
		outpointStr := hex.EncodeToString(key.OutPoint().RawData())
		claimedCotas = append(claimedCotas, biz.ClaimedCotaNftKvPair{
			BlockNumber: blockNumber,
			CotaId:      hex.EncodeToString(key.NftId().CotaId().RawData()),
			CotaIdCRC:   crc32.ChecksumIEEE([]byte(cotaId)),
			TokenIndex:  binary.BigEndian.Uint32(key.NftId().Index().RawData()),
			OutPoint:    outpointStr,
			OutPointCrc: crc32.ChecksumIEEE([]byte(outpointStr)),
			LockHash:    lockHashStr,
			LockHashCrc: lockHashCRC32,
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
			return claimedCotas, withdrawCotas, err
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
	return claimedCotas, withdrawCotas, nil
}

func generateTransferV1ToV2KvPairs(blockNumber uint64, entry biz.Entry, rp transferCotaKvPairRepo) ([]biz.ClaimedCotaNftKvPair, []biz.WithdrawCotaNftKvPair, error) {
	var (
		claimedCotaKeyVec *smt.ClaimCotaNFTKeyVec
		withdrawKeyVec    *smt.WithdrawalCotaNFTKeyV1Vec
		withdrawValueVec  *smt.WithdrawalCotaNFTValueV1Vec
		claimedCotas      []biz.ClaimedCotaNftKvPair
		withdrawCotas     []biz.WithdrawCotaNftKvPair
	)
	if entry.Version == 1 {
		entries, err := smt.TransferCotaNFTV1EntriesFromSlice(entry.InputType[1:], false)
		if err != nil {
			return claimedCotas, withdrawCotas, err
		}
		claimedCotaKeyVec = entries.ClaimKeys()
		withdrawKeyVec = entries.WithdrawalKeys()
		withdrawValueVec = entries.WithdrawalValues()
	} else {
		entries, err := smt.TransferCotaNFTV2EntriesFromSlice(entry.InputType[1:], false)
		claimedCotaKeyVec = entries.ClaimKeys()
		withdrawKeyVec = entries.WithdrawalKeys()
		withdrawValueVec = entries.WithdrawalValues()
		if err != nil {
			return claimedCotas, withdrawCotas, err
		}
	}
	senderLock, lockHashStr, lockHashCRC32, err := GenerateSenderLock(entry)
	if err != nil {
		return claimedCotas, withdrawCotas, err
	}
	if err := rp.FindOrCreateScript(context.TODO(), &senderLock); err != nil {
		return claimedCotas, withdrawCotas, err
	}
	for i := uint(0); i < claimedCotaKeyVec.Len(); i++ {
		key := claimedCotaKeyVec.Get(i)
		cotaId := hex.EncodeToString(key.NftId().CotaId().RawData())
		outpointStr := hex.EncodeToString(key.OutPoint().RawData())
		claimedCotas = append(claimedCotas, biz.ClaimedCotaNftKvPair{
			BlockNumber: blockNumber,
			CotaId:      hex.EncodeToString(key.NftId().CotaId().RawData()),
			CotaIdCRC:   crc32.ChecksumIEEE([]byte(cotaId)),
			TokenIndex:  binary.BigEndian.Uint32(key.NftId().Index().RawData()),
			OutPoint:    outpointStr,
			OutPointCrc: crc32.ChecksumIEEE([]byte(outpointStr)),
			LockHash:    lockHashStr,
			LockHashCrc: lockHashCRC32,
		})
	}
	for i := uint(0); i < withdrawKeyVec.Len(); i++ {
		key := withdrawKeyVec.Get(i)
		value := withdrawValueVec.Get(i)
		cotaId := hex.EncodeToString(key.NftId().CotaId().RawData())
		outpointStr := hex.EncodeToString(key.OutPoint().RawData())
		receiverLock := GenerateReceiverLock(value.ToLock().RawData())
		err = rp.FindOrCreateScript(context.TODO(), &receiverLock)
		if err != nil {
			return claimedCotas, withdrawCotas, err
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
	return claimedCotas, withdrawCotas, nil
}
