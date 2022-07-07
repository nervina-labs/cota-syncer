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

var _ biz.ClaimedCotaNftKvPairRepo = (*claimedCotaNftKvPairRepo)(nil)

type ClaimedCotaNftKvPair struct {
	ID          uint `gorm:"primaryKey"`
	BlockNumber uint64
	CotaId      string
	CotaIdCRC   uint32
	TokenIndex  uint32
	OutPoint    string
	OutPointCrc uint32
	LockHash    string
	LockHashCrc uint32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type claimedCotaNftKvPairRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewClaimedCotaNftKvPairRepo(data *Data, logger *logger.Logger) biz.ClaimedCotaNftKvPairRepo {
	return &claimedCotaNftKvPairRepo{
		data:   data,
		logger: logger,
	}
}

func (rp claimedCotaNftKvPairRepo) CreateClaimedCotaNftKvPair(ctx context.Context, c *biz.ClaimedCotaNftKvPair) error {
	if err := rp.data.db.WithContext(ctx).Create(c).Error; err != nil {
		return err
	}
	return nil
}

func (rp claimedCotaNftKvPairRepo) DeleteClaimedCotaNftKvPairs(ctx context.Context, blockNumber uint64) error {
	if err := rp.data.db.WithContext(ctx).Where("block_number = ?", blockNumber).Delete(ClaimedCotaNftKvPair{}).Error; err != nil {
		return err
	}
	return nil
}

func (rp claimedCotaNftKvPairRepo) ParseClaimedCotaEntries(blockNumber uint64, entry biz.Entry) ([]biz.HoldCotaNftKvPair, []biz.ClaimedCotaNftKvPair, error) {
	var (
		holdCotaKeyVec    *smt.HoldCotaNFTKeyVec
		holdCotaValueVec  *smt.HoldCotaNFTValueVec
		claimedCotaKeyVec *smt.ClaimCotaNFTKeyVec
		holdCotas         []biz.HoldCotaNftKvPair
		claimedCotas      []biz.ClaimedCotaNftKvPair
	)

	if entry.Version == 0 {
		entries, err := smt.ClaimCotaNFTEntriesFromSlice(entry.InputType[1:], false)
		if err != nil {
			return holdCotas, claimedCotas, err
		}
		holdCotaKeyVec = entries.HoldKeys()
		holdCotaValueVec = entries.HoldValues()
		claimedCotaKeyVec = entries.ClaimKeys()
	} else {
		entries, err := smt.ClaimCotaNFTV2EntriesFromSlice(entry.InputType[1:], false)
		if err != nil {
			return holdCotas, claimedCotas, err
		}
		holdCotaKeyVec = entries.HoldKeys()
		holdCotaValueVec = entries.HoldValues()
		claimedCotaKeyVec = entries.ClaimKeys()
	}
	lockHash, err := entry.LockScript.Hash()
	if err != nil {
		return holdCotas, claimedCotas, err
	}
	lockHashStr := lockHash.String()[2:]
	lockHashCRC32 := crc32.ChecksumIEEE([]byte(lockHashStr))
	for i := uint(0); i < holdCotaKeyVec.Len(); i++ {
		key := holdCotaKeyVec.Get(i)
		value := holdCotaValueVec.Get(i)
		holdCotas = append(holdCotas, biz.HoldCotaNftKvPair{
			BlockNumber:    blockNumber,
			CotaId:         hex.EncodeToString(key.CotaId().RawData()),
			TokenIndex:     binary.BigEndian.Uint32(key.Index().RawData()),
			State:          value.State().AsSlice()[0],
			Configure:      value.Configure().AsSlice()[0],
			Characteristic: hex.EncodeToString(value.Characteristic().RawData()),
			LockHash:       lockHashStr,
			LockHashCRC:    lockHashCRC32,
		})
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
	return holdCotas, claimedCotas, nil
}

func (rp claimedCotaNftKvPairRepo) ParseClaimedUpdateCotaEntries(blockNumber uint64, entry biz.Entry) ([]biz.HoldCotaNftKvPair, []biz.ClaimedCotaNftKvPair, error) {
	var (
		holdCotaKeyVec    *smt.HoldCotaNFTKeyVec
		holdCotaValueVec  *smt.HoldCotaNFTValueVec
		claimedCotaKeyVec *smt.ClaimCotaNFTKeyVec
		holdCotas         []biz.HoldCotaNftKvPair
		claimedCotas      []biz.ClaimedCotaNftKvPair
	)
	if entry.Version == 0 {
		entries, err := smt.ClaimUpdateCotaNFTEntriesFromSlice(entry.InputType[1:], false)
		if err != nil {
			return holdCotas, claimedCotas, err
		}
		holdCotaKeyVec = entries.HoldKeys()
		holdCotaValueVec = entries.HoldValues()
		claimedCotaKeyVec = entries.ClaimKeys()
	} else {
		entries, err := smt.ClaimUpdateCotaNFTV2EntriesFromSlice(entry.InputType[1:], false)
		if err != nil {
			return holdCotas, claimedCotas, err
		}
		holdCotaKeyVec = entries.HoldKeys()
		holdCotaValueVec = entries.HoldValues()
		claimedCotaKeyVec = entries.ClaimKeys()
	}

	lockHash, err := entry.LockScript.Hash()
	if err != nil {
		return holdCotas, claimedCotas, err
	}
	lockHashStr := lockHash.String()[2:]
	lockHashCRC32 := crc32.ChecksumIEEE([]byte(lockHashStr))
	for i := uint(0); i < holdCotaKeyVec.Len(); i++ {
		key := holdCotaKeyVec.Get(i)
		value := holdCotaValueVec.Get(i)
		holdCotas = append(holdCotas, biz.HoldCotaNftKvPair{
			BlockNumber:    blockNumber,
			CotaId:         hex.EncodeToString(key.CotaId().RawData()),
			TokenIndex:     binary.BigEndian.Uint32(key.Index().RawData()),
			State:          value.State().AsSlice()[0],
			Configure:      value.Configure().AsSlice()[0],
			Characteristic: hex.EncodeToString(value.Characteristic().RawData()),
			LockHash:       lockHashStr,
			LockHashCRC:    lockHashCRC32,
		})
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
	return holdCotas, claimedCotas, nil
}
