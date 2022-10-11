package service

import (
	"context"
	"encoding/hex"

	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/data"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
)

var _ Service = (*RegisterLockService)(nil)

type RegisterLockService struct {
	lockScriptUsecase *biz.RegisterLockScriptUsecase
	logger            *logger.Logger
	client            *data.CkbNodeClient
}

func NewRegisterLockService(lockScriptUsecase *biz.RegisterLockScriptUsecase, logger *logger.Logger, client *data.CkbNodeClient) *RegisterLockService {
	return &RegisterLockService{
		lockScriptUsecase: lockScriptUsecase,
		logger:            logger,
		client:            client,
	}
}

func (s RegisterLockService) Start(ctx context.Context, _ string) error {
	s.logger.Info(ctx, "register lock script service started")
	page := 0
	for {
		queryInfos, err := s.lockScriptUsecase.FindRegisterQueryInfos(ctx, page, pageSize)
		if err != nil {
			return err
		}
		if len(queryInfos) == 0 {
			break
		}
		page += 1

		if err = s.parseLockScripts(ctx, queryInfos); err != nil {
			return err
		}
	}
	return nil
}

func (s RegisterLockService) Stop(ctx context.Context) error {
	s.logger.Info(ctx, "register lock script service stopped")

	return nil
}

func (s RegisterLockService) parseLockScripts(ctx context.Context, infos []biz.RegisterQueryInfo) (err error) {
	var block *ckbTypes.Block
	for _, info := range infos {
		block, err = s.client.Rpc.GetBlockByNumber(ctx, info.BlockNumber)
		if err != nil {
			return
		}
		for _, tx := range block.Transactions {
			if len(tx.Outputs) < 2 {
				continue
			}
			if err = s.parseLockScript(ctx, tx, info); err != nil {
				return
			}
		}
	}
	return
}

func (s RegisterLockService) parseLockScript(ctx context.Context, tx *ckbTypes.Transaction, info biz.RegisterQueryInfo) error {
	// Find CoTA registry transactions with lockHash
	for _, output := range tx.Outputs {
		lockHash, err := output.Lock.Hash()
		if err != nil {
			return err
		}
		if lockHash.String()[2:] != info.LockHash {
			continue
		}
		hashType, err := output.Lock.HashType.Serialize()
		if err != nil {
			return err
		}
		lock := biz.Script{
			CodeHash: hex.EncodeToString(output.Lock.CodeHash[:]),
			HashType: hex.EncodeToString(hashType),
			Args:     hex.EncodeToString(output.Lock.Args),
		}
		if err = s.lockScriptUsecase.FindOrCreateScript(ctx, &lock); err != nil {
			return err
		}
		if err = s.lockScriptUsecase.CreateRegisterLock(ctx, info.LockHash, lock.ID); err != nil {
			return err
		}
	}
	return nil
}
