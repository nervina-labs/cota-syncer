/*
Copyright (c) 2020 go-kratos, All Rights Reserved
Distributed under MIT license. See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

package app

import (
	"context"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/data"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/logger"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/service"
	"os"
	"time"
)

type Option func(o *options)

type options struct {
	id      string
	name    string
	version string

	ctx  context.Context
	sigs []os.Signal

	logger      *logger.Logger
	stopTimeout time.Duration
	services    []service.Service
	migration   *data.DBMigration
}

func ID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

func Name(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func Version(version string) Option {
	return func(o *options) {
		o.version = version
	}
}

func Context(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

func Logger(logger *logger.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func Signal(sigs ...os.Signal) Option {
	return func(o *options) {
		o.sigs = sigs
	}
}

func StopTimeout(t time.Duration) Option {
	return func(o *options) {
		o.stopTimeout = t
	}
}

func Services(srvs ...service.Service) Option {
	return func(o *options) {
		o.services = srvs
	}
}

func Migration(m *data.DBMigration) Option {
	return func(o *options) {
		o.migration = m
	}
}
