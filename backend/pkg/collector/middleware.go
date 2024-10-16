package collector

import (
	"context"

	"backend/pkg/config"
	"backend/pkg/product"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"
)

var (
	KeyRequestCollector kit.Key = kit.KeyBase + "request:collector"
)

func RequestCollector(ctx context.Context) *Collector {
	return ctx.Value(KeyRequestCollector).(*Collector) // nolint:forcetypeassert,errcheck
}

type CollectorMiddleware struct {
	config              config.Config
	observer            *kit.Observer
	collectorRepository *CollectorRepository
}

func NewCollectorMiddleware(observer *kit.Observer, collectorRepository *CollectorRepository,
	config config.Config) *CollectorMiddleware {
	return &CollectorMiddleware{
		config:              config,
		observer:            observer,
		collectorRepository: collectorRepository,
	}
}

func (self *CollectorMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestCtx := ctx.Request().Context()
		requestProduct := product.RequestProduct(requestCtx)

		collector, err := self.collectorRepository.GetByID(requestCtx, ctx.Param("collector_id"))
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		if collector == nil {
			return kit.HTTPErrInvalidRequest
		}

		if collector.ProductID != requestProduct.ID {
			return kit.HTTPErrUnauthorized
		}

		if collector.DeletedAt != nil {
			return kit.HTTPErrUnauthorized
		}

		ctx.SetRequest(ctx.Request().WithContext(context.WithValue(requestCtx, KeyRequestCollector, collector)))

		return next(ctx)
	}
}
