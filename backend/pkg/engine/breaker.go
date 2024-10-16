package engine

import (
	"context"
	"sync"
	"time"

	"backend/pkg/config"

	"github.com/neoxelox/errors"
	"github.com/neoxelox/kit"
)

const (
	ENGINE_BREAKER_KEY          = "engine:breaker"
	ENGINE_BREAKER_TIMEOUT      = 30 * time.Second
	ENGINE_BREAKER_MIN_FAILURES = 25
)

var (
	ErrEngineBreakerGeneric = errors.New("engine breaker failed")
	ErrEngineBreakerOpened  = errors.New("engine breaker opened")
)

type EngineBreaker struct {
	config   config.Config
	observer *kit.Observer
	cache    *kit.Cache
	mutex    sync.Mutex
}

func NewEngineBreaker(observer *kit.Observer, cache *kit.Cache, config config.Config) *EngineBreaker {
	return &EngineBreaker{
		config:   config,
		observer: observer,
		cache:    cache,
		mutex:    sync.Mutex{},
	}
}

func (self *EngineBreaker) IsOpen(ctx context.Context) bool {
	var failures int

	err := self.cache.Get(ctx, ENGINE_BREAKER_KEY, &failures)
	if err != nil {
		if kit.ErrCacheMiss.Is(err) {
			return false
		}

		self.observer.Error(ctx, ErrEngineBreakerGeneric.Raise().Cause(err))

		return true
	}

	if failures >= ENGINE_BREAKER_MIN_FAILURES {
		return true
	}

	return false
}

// TODO: Refactor open functionality with Redis transactional pipelines!
func (self *EngineBreaker) open(ctx context.Context, failures int, timeout time.Duration) error {
	err := self.cache.Set(ctx, ENGINE_BREAKER_KEY, failures, &timeout)
	if err != nil {
		return ErrEngineBreakerGeneric.Raise().Cause(err)
	}

	if failures == ENGINE_BREAKER_MIN_FAILURES {
		self.observer.Error(ctx, ErrEngineBreakerOpened.Raise().
			Extra(map[string]any{"timeout": int(timeout.Seconds())}))
	}

	return nil
}

func (self *EngineBreaker) Open(ctx context.Context) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	var failures int

	err := self.cache.Get(ctx, ENGINE_BREAKER_KEY, &failures)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return ErrEngineBreakerGeneric.Raise().Cause(err)
	}

	failures++

	return self.open(ctx, failures, ENGINE_BREAKER_TIMEOUT)
}

func (self *EngineBreaker) Force(ctx context.Context, timeout time.Duration) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	return self.open(ctx, ENGINE_BREAKER_MIN_FAILURES, timeout)
}

func (self *EngineBreaker) Close(ctx context.Context) error {
	err := self.cache.Delete(ctx, ENGINE_BREAKER_KEY)
	if err != nil {
		return ErrEngineBreakerGeneric.Raise().Cause(err)
	}

	self.observer.Info(ctx, "Engine circuit breaker closed")

	return nil
}
