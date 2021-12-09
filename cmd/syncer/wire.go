//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/app"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/config"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/data"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/logger"
	"github.com/nervina-labs/compact-nft-entries-syncer/internal/service"
)

func initApp(*config.Database, *config.CkbNode, *logger.Logger) (*app.App, func(), error) {
	panic(wire.Build(data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
