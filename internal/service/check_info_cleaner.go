package service

import (
	"context"
	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/data"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	"golang.org/x/sync/errgroup"
	"time"
)

var _ Service = (*CheckInfoCleanerService)(nil)

type CheckInfoCleanerService struct {
	checkInfoUsecase *biz.CheckInfoUsecase
	logger           *logger.Logger
	client           *data.CkbNodeClient
}

func NewCheckInfoService(checkInfoUsecase *biz.CheckInfoUsecase, logger *logger.Logger, client *data.CkbNodeClient) *CheckInfoCleanerService {
	return &CheckInfoCleanerService{
		checkInfoUsecase: checkInfoUsecase,
		logger:           logger,
		client:           client,
	}
}

func (scv CheckInfoCleanerService) clean(ctx context.Context, checkType biz.CheckType) error {
	return scv.checkInfoUsecase.Clean(ctx, checkType)
}

func (scv CheckInfoCleanerService) Start(ctx context.Context, mode string) error {
	scv.logger.Info(ctx, "Successfully started the check info cleaner~")

	var interval time.Duration

	if mode == "normal" {
		interval = 30 * time.Minute
	} else {
		interval = 1 * time.Minute
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				scv.logger.Infof(ctx, "cleaner received cancel signal %v", ctx.Err())
			case <-time.After(interval):
				eg, ctx := errgroup.WithContext(ctx)
				checkTypes := []biz.CheckType{biz.SyncBlock, biz.SyncMetadata}
				for _, checkType := range checkTypes {
					cType := checkType
					eg.Go(func() error {
						return scv.clean(ctx, cType)
					})
				}
				if err := eg.Wait(); err != nil {
					scv.logger.Errorf(ctx, "clean check info failed, %v", err)
				}
			}
		}
	}()
	return nil
}

func (scv CheckInfoCleanerService) Stop(ctx context.Context) error {
	scv.client.Rpc.Close()
	scv.logger.Infof(ctx, "Successfully closed the cleaner service~")

	return nil
}
