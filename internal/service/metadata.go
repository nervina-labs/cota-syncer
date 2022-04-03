package service

import (
	"context"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/data"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
	"time"
)

type MetadataSyncService struct {
	checkInfoUsecase *biz.CheckInfoUsecase
	logger           *logger.Logger
	client           *data.CkbNodeClient
	status           chan struct{}
	systemScripts    data.SystemScripts
	metadataSyncer   data.MetadataSyncer
}

func NewMetadataSyncService(checkInfoUsecase *biz.CheckInfoUsecase, logger *logger.Logger, client *data.CkbNodeClient, systemScripts data.SystemScripts, metadataSyncer data.MetadataSyncer) *MetadataSyncService {
	return &MetadataSyncService{
		checkInfoUsecase: checkInfoUsecase,
		logger:           logger,
		client:           client,
		status:           make(chan struct{}, 1),
		systemScripts:    systemScripts,
		metadataSyncer:   metadataSyncer,
	}
}

func (s *MetadataSyncService) Start(ctx context.Context, mode string) error {
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
				if mode == "normal" {
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()
	return nil
}

func (s *MetadataSyncService) Stop(ctx context.Context) error {
	s.client.Rpc.Close()
	for {
		select {
		case <-s.status:
			s.logger.Info(ctx, "Successfully closed the metadata sync service~")
			return nil
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (s *MetadataSyncService) sync(ctx context.Context) {
	checkInfo := biz.CheckInfo{CheckType: biz.SyncMetadata}
	err := s.checkInfoUsecase.LastCheckInfo(ctx, &checkInfo)
	if err != nil {
		s.logger.Errorf(ctx, "get %s check info error: %v", checkInfo.CheckType.String(), err)
	}
	tipBlockNumber, err := s.client.Rpc.GetTipBlockNumber(ctx)
	if err != nil {
		s.logger.Errorf(ctx, "get tip block number rpc error: %v", err)
	}
	s.logger.Infof(ctx, "check tip block number: %v, tip block number: %v", checkInfo.BlockNumber, tipBlockNumber)
	if checkInfo.BlockNumber > tipBlockNumber {
		return
	}
	targetBlockNumber := checkInfo.BlockNumber + 1
	if targetBlockNumber > tipBlockNumber {
		return
	}
	targetBlock, err := s.client.Rpc.GetBlockByNumber(ctx, targetBlockNumber)
	if err != nil {
		s.logger.Errorf(ctx, "get block %d rpc error: %v", targetBlockNumber, err)
		return
	}
	// rollback
	if isForked(checkInfo, targetBlock) {
		s.logger.Info(ctx, "forked")
		err = s.rollback(ctx, checkInfo.BlockNumber)
		if err != nil {
			s.logger.Errorf(ctx, "rollback %s error: %v", checkInfo.CheckType.String(), err)
		}
		return
	}
	// save key pairs
	checkInfo.BlockNumber = targetBlockNumber
	checkInfo.BlockHash = targetBlock.Header.Hash.String()[2:]
	err = s.syncMetadata(ctx, targetBlock, checkInfo)
	if err != nil {
		s.logger.Errorf(ctx, "save %s kv pairs error: %v", checkInfo.CheckType.String(), err)
	}
}

func (s *MetadataSyncService) syncMetadata(ctx context.Context, block *ckbTypes.Block, checkInfo biz.CheckInfo) error {
	return s.metadataSyncer.Sync(ctx, block, checkInfo, s.systemScripts)
}

func (s *MetadataSyncService) rollback(ctx context.Context, blockNumber uint64) error {
	return s.metadataSyncer.Rollback(ctx, blockNumber)
}
