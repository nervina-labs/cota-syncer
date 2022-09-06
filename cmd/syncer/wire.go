//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/nervina-labs/cota-syncer/internal/app"
	"github.com/nervina-labs/cota-syncer/internal/biz"
	"github.com/nervina-labs/cota-syncer/internal/config"
	"github.com/nervina-labs/cota-syncer/internal/data"
	"github.com/nervina-labs/cota-syncer/internal/logger"
	"github.com/nervina-labs/cota-syncer/internal/service"
)

func initApp(*config.Database, *config.CkbNode, *logger.Logger) (*app.App, func(), error) {
	panic(wire.Build(data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
