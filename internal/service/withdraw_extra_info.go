package service

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"

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
	var lock biz.Script
	var index uint32
	for _, v := range queryInfos {
		block, err = s.client.Rpc.GetBlockByNumber(ctx, v.BlockNumber)
		if err != nil {
			return err
		}
		for _, tx := range block.Transactions {
			if hex.EncodeToString(tx.Hash[12:]) == v.OutPoint[:20] {
				index = binary.LittleEndian.Uint32(([]byte)(v.OutPoint[20:]))
				if index >= (uint32)(len(tx.Outputs)) {
					return errors.New("out_point index error")
				}
				lockScript := tx.Outputs[index].Lock
				lock = biz.Script{
					CodeHash: hex.EncodeToString(lockScript.CodeHash[:]),
					HashType: (string)(lockScript.HashType),
					Args:     hex.EncodeToString(lockScript.Args),
				}
				if err = s.extraInfoUsecase.FindOrCreateScript(ctx, &lock); err != nil {
					return err
				}
				if err = s.extraInfoUsecase.CreateExtraInfo(ctx, v.OutPoint, tx.Hash.String()[2:], lock.ID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s WithdrawExtraInfoService) Stop(ctx context.Context) error {
	s.logger.Info(ctx, "withdraw extra info service stopped")

	return nil
}
