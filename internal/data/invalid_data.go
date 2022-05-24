package data

import (
	"context"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
)

var _ biz.InvalidDataRepo = (*invalidDataRepo)(nil)

type invalidDataRepo struct {
	data   *Data
	logger *logger.Logger
}

func NewInvalidDateRepo(data *Data, logger *logger.Logger) biz.InvalidDataRepo {
	return &invalidDataRepo{
		data:   data,
		logger: logger,
	}
}

func (rp invalidDataRepo) Clean(ctx context.Context, blockNumber uint64) error {
	define := DefineCotaNftKvPair{
		Total:     0,
		Issued:    0,
		Configure: 0,
	}
	if err := rp.data.db.Debug().WithContext(ctx).Where("block_number < ? and total = ? and issued = ? and configure = ?", blockNumber, define.Total, define.Issued, define.Configure).Delete(DefineCotaNftKvPair{}).Error; err != nil {
		return err
	}

	hold := HoldCotaNftKvPair{
		State:          0,
		Configure:      0,
		Characteristic: "0000000000000000000000000000000000000000",
	}
	if err := rp.data.db.Debug().WithContext(ctx).Where("block_number < ? and state = ? and configure = ? and characteristic = ?", blockNumber, hold.State, hold.Configure, hold.Characteristic).Delete(HoldCotaNftKvPair{}).Error; err != nil {
		return err
	}

	return nil
}
