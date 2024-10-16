package exporter

import (
	"context"

	"backend/pkg/config"
	"backend/pkg/product"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"
)

var (
	KeyRequestExporter kit.Key = kit.KeyBase + "request:exporter"
)

func RequestExporter(ctx context.Context) *Exporter {
	return ctx.Value(KeyRequestExporter).(*Exporter) // nolint:forcetypeassert,errcheck
}

type ExporterMiddleware struct {
	config             config.Config
	observer           *kit.Observer
	exporterRepository *ExporterRepository
}

func NewExporterMiddleware(observer *kit.Observer, exporterRepository *ExporterRepository,
	config config.Config) *ExporterMiddleware {
	return &ExporterMiddleware{
		config:             config,
		observer:           observer,
		exporterRepository: exporterRepository,
	}
}

func (self *ExporterMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestCtx := ctx.Request().Context()
		requestProduct := product.RequestProduct(requestCtx)

		exporter, err := self.exporterRepository.GetByID(requestCtx, ctx.Param("exporter_id"))
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		if exporter == nil {
			return kit.HTTPErrInvalidRequest
		}

		if exporter.ProductID != requestProduct.ID {
			return kit.HTTPErrUnauthorized
		}

		if exporter.DeletedAt != nil {
			return kit.HTTPErrUnauthorized
		}

		ctx.SetRequest(ctx.Request().WithContext(context.WithValue(requestCtx, KeyRequestExporter, exporter)))

		return next(ctx)
	}
}
