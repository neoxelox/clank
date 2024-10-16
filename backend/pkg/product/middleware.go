package product

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
	"backend/pkg/organization"
)

var (
	KeyRequestProduct kit.Key = kit.KeyBase + "request:product"
)

func RequestProduct(ctx context.Context) *Product {
	return ctx.Value(KeyRequestProduct).(*Product) // nolint:forcetypeassert,errcheck
}

type ProductMiddleware struct {
	config            config.Config
	observer          *kit.Observer
	productRepository *ProductRepository
}

func NewProductMiddleware(observer *kit.Observer, productRepository *ProductRepository,
	config config.Config) *ProductMiddleware {
	return &ProductMiddleware{
		config:            config,
		observer:          observer,
		productRepository: productRepository,
	}
}

func (self *ProductMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestCtx := ctx.Request().Context()
		requestOrganization := organization.RequestOrganization(requestCtx)

		product, err := self.productRepository.GetByID(requestCtx, ctx.Param("product_id"))
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		if product == nil {
			return kit.HTTPErrInvalidRequest
		}

		if product.OrganizationID != requestOrganization.ID {
			return kit.HTTPErrUnauthorized
		}

		if product.DeletedAt != nil {
			return kit.HTTPErrUnauthorized
		}

		ctx.SetRequest(ctx.Request().WithContext(context.WithValue(requestCtx, KeyRequestProduct, product)))

		return next(ctx)
	}
}
