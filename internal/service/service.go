package service

import (
	"context"
	"github.com/google/wire"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/data"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/logger"
	"time"
)

var ProviderSet = wire.NewSet(NewSyncService)

type SyncService struct {
	tipBlock *biz.TipBlockUsecase
	logger   *logger.Logger
	client   *data.CkbNodeClient
}

func (s *SyncService) Start(ctx context.Context) error {
	s.logger.Info(ctx, "Successfully started the sync service~")
	for {
		number, err := s.client.Rpc.GetTipBlockNumber(ctx)
		//block, err := s.client.Rpc.GetBlockByNumber(ctx, 30000)

		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
		return s.sync(ctx, number)
	}
}

func (s *SyncService) sync(ctx context.Context, blockNumber uint64) error {
	s.logger.Infof(ctx,"current block number: %v", blockNumber)
	return nil
}

func (s *SyncService) Stop(ctx context.Context) error {
	s.logger.Info(ctx,"Successfully closed the sync service~")
	return nil
}

func NewSyncService(tipBlock *biz.TipBlockUsecase, logger *logger.Logger, client *data.CkbNodeClient) *SyncService {
	return &SyncService{
		tipBlock: tipBlock,
		logger:   logger,
		client:   client,
	}
}

type Service interface {
	Start(context.Context) error
	Stop(context.Context) error
}
