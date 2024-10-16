package util

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
)

var (
	KeyRateLimit kit.Key = kit.KeyBase + "ratelimit:"
)

type RateLimitMiddleware struct {
	config   config.Config
	observer *kit.Observer
	limiter  *kit.Limiter
}

func NewRateLimitMiddleware(observer *kit.Observer, limiter *kit.Limiter, config config.Config) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		config:   config,
		observer: observer,
		limiter:  limiter,
	}
}

func (self *RateLimitMiddleware) Handle(limit int, period time.Duration) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			requestCtx := ctx.Request().Context()

			remaining, err := self.limiter.Limit(requestCtx,
				string(KeyRateLimit)+ctx.Path()+ctx.RealIP(), limit, period)
			if err != nil {
				return kit.HTTPErrServerGeneric.Cause(err)
			}

			if remaining < 0 {
				return kit.HTTPErrRateLimited
			}

			return next(ctx)
		}
	}
}
