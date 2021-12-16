package service

import (
	"context"
	"github.com/google/wire"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/data"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
	"time"
)

var ProviderSet = wire.NewSet(NewSyncService)

type SyncService struct {
	checkInfoUsecase *biz.CheckInfoUsecase
	kvPairUsecase    *biz.SyncKvPairUsecase
	logger           *logger.Logger
	client           *data.CkbNodeClient
	status           chan struct{}
	systemScripts    data.SystemScripts
	blockParser      data.BlockParser
}

func (s *SyncService) Start(ctx context.Context) error {
	s.logger.Info(ctx, "Successfully started the sync service~")
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.status <- struct{}{}
				s.logger.Infof(ctx, "receive cancel signal %v", ctx.Err())
				return
			default:
				s.sync(ctx)
				time.Sleep(1 * time.Second)
			}
		}
	}()
	return nil
}

func (s *SyncService) sync(ctx context.Context) {
	tipBlockNumber, err := s.client.Rpc.GetTipBlockNumber(ctx)
	if err != nil {
		s.logger.Errorf(ctx, "get tip block number rpc error: %v", err)
	}
	s.logger.Infof(ctx, "block_number: %v", tipBlockNumber)
	checkInfo := biz.CheckInfo{CheckType: biz.SyncEvent}
	err = s.checkInfoUsecase.FindOrCreate(ctx, &checkInfo)
	if err != nil {
		s.logger.Errorf(ctx, "get check info error: %v", err)
	}
	if checkInfo.BlockNumber > tipBlockNumber {
		return
	}

	tipBlock, err := s.client.Rpc.GetBlockByNumber(ctx, checkInfo.BlockNumber)
	// rollback
	if checkInfo.BlockHash != tipBlock.Header.ParentHash.String() {
		s.logger.Info(ctx, "forked")
		err = s.rollback(ctx, checkInfo.BlockNumber)
		if err != nil {
			s.logger.Errorf(ctx, "rollback error: %v", err)
		}
		return
	}
	// save key pairs
	err = s.saveKvPairs(ctx, tipBlock)
	if err != nil {
		s.logger.Errorf(ctx, "save kv pairs error: %v", err)
	}
}

func (s *SyncService) saveKvPairs(ctx context.Context, block *ckbTypes.Block) error {
	kvPair, err := s.blockParser.Parse(block, s.systemScripts)
	if err != nil {
		return err
	}
	return s.kvPairUsecase.CreateKvPairs(ctx, &kvPair)
}

func (s *SyncService) rollback(ctx context.Context, blockNumber uint64) error {
	return s.kvPairUsecase.DeleteKvPairs(ctx, blockNumber)
}

func (s *SyncService) Stop(ctx context.Context) error {
	s.client.Rpc.Close()
	for {
		select {
		case <-s.status:
			s.logger.Info(ctx, "Successfully closed the sync service~")
			return nil
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func NewSyncService(checkInfoUsecase *biz.CheckInfoUsecase, kvPairUsecase *biz.SyncKvPairUsecase, logger *logger.Logger, client *data.CkbNodeClient, systemScripts data.SystemScripts, blockParser data.BlockParser) *SyncService {
	return &SyncService{
		checkInfoUsecase: checkInfoUsecase,
		kvPairUsecase:    kvPairUsecase,
		logger:           logger,
		client:           client,
		status:           make(chan struct{}, 1),
		systemScripts:    systemScripts,
		blockParser:      blockParser,
	}
}

type Service interface {
	Start(context.Context) error
	Stop(context.Context) error
}
