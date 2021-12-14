/*
Copyright (c) 2020 go-kratos, All Rights Reserved
Distributed under MIT license. See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

package app

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type AppInfo interface {
	ID() string
	Name() string
	Version() string
}

type App struct {
	options options
	ctx     context.Context
	cancel  func()
	lk      sync.Mutex
}

func NewApp(opts ...Option) *App {
	o := options{
		ctx:         context.Background(),
		sigs:        []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		stopTimeout: 10 * time.Second,
	}
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}
	for _, opt := range opts {
		opt(&o)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &App{
		options: o,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// ID returns app instance id.
func (a *App) ID() string { return a.options.id }

// Name returns services name.
func (a *App) Name() string { return a.options.name }

// Version returns app version.
func (a *App) Version() string { return a.options.version }

func (a *App) Run() error {
	ctx := NewContext(a.ctx, a)
	eg, ctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	if err := a.options.migration.Up(); err != nil {
		a.options.logger.Errorf(context.TODO(), "DB Migration failed: %v", err)
		return err
	}
	for _, srv := range a.options.services {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done()
			sctx, cancel := context.WithTimeout(NewContext(context.Background(), a), a.options.stopTimeout)
			defer cancel()
			return srv.Stop(sctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start(ctx)
		})
	}
	wg.Wait()
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.options.sigs...)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				err := a.Stop()
				if err != nil {
					a.options.logger.Errorf(context.TODO(), "failed to stop app: %v, %v", a.Name(), err)
					return err
				}
			}
		}
	})
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
		a.options.logger.Infof(context.TODO(), "successfully stop the app: %v", a.Name())
	}
	return nil
}

type appKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, s AppInfo) context.Context {
	return context.WithValue(ctx, appKey{}, s)
}
