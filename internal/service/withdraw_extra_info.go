package service

import (
	"context"
	"encoding/hex"

	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/data"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
)

var _ Service = (*WithdrawExtraInfoService)(nil)

type WithdrawExtraInfoService struct {
	extraInfoUsecase *biz.WithdrawExtraInfoUsecase
	logger           *logger.Logger
	client           *data.CkbNodeClient
}

func NewWithdrawExtraInfoService(extraInfoUsecase *biz.WithdrawExtraInfoUsecase, logger *logger.Logger, client *data.CkbNodeClient) *WithdrawExtraInfoService {
	return &WithdrawExtraInfoService{
		extraInfoUsecase: extraInfoUsecase,
		logger:           logger,
		client:           client,
	}
}

func (s WithdrawExtraInfoService) Start(ctx context.Context, _ string) error {
	s.logger.Info(ctx, "withdraw extra info service started")
	queryInfos, err := s.extraInfoUsecase.FindAllQueryInfos(ctx)
	if err != nil {
		return err
	}
	var block *ckbTypes.Block
	for _, info := range queryInfos {
		block, err = s.client.Rpc.GetBlockByNumber(ctx, info.BlockNumber)
		if err != nil {
			return err
		}
		for _, tx := range block.Transactions {
			if err = s.parseExtraInfo(ctx, tx, info); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s WithdrawExtraInfoService) Stop(ctx context.Context) error {
	s.logger.Info(ctx, "withdraw extra info service stopped")

	return nil
}

func (s WithdrawExtraInfoService) parseExtraInfo(ctx context.Context, tx *ckbTypes.Transaction, info biz.WithdrawQueryInfo) error {
	for _, input := range tx.Inputs {
		// Find CoTA transactions with txHash and outPoint
		if hex.EncodeToString(input.PreviousOutput.TxHash[12:]) != info.OutPoint[:40] {
			continue
		}
		for _, output := range tx.Outputs {
			lockHash, err := output.Lock.Hash()
			if err != nil {
				return err
			}
			if lockHash.String()[2:] == info.LockHash {
				hashType, err := output.Lock.HashType.Serialize()
				if err != nil {
					return err
				}
				lock := biz.Script{
					CodeHash: hex.EncodeToString(output.Lock.CodeHash[:]),
					HashType: hex.EncodeToString(hashType),
					Args:     hex.EncodeToString(output.Lock.Args),
				}
				if err = s.extraInfoUsecase.FindOrCreateScript(ctx, &lock); err != nil {
					return err
				}
				if err = s.extraInfoUsecase.CreateExtraInfo(ctx, info.OutPoint, tx.Hash.String()[2:], lock.ID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
