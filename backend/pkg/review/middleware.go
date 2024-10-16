package review

import (
	"context"

	"backend/pkg/config"
	"backend/pkg/product"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"
)

var (
	KeyRequestReview kit.Key = kit.KeyBase + "request:review"
)

func RequestReview(ctx context.Context) *Review {
	return ctx.Value(KeyRequestReview).(*Review) // nolint:forcetypeassert,errcheck
}

type ReviewMiddleware struct {
	config           config.Config
	observer         *kit.Observer
	reviewRepository *ReviewRepository
}

func NewReviewMiddleware(observer *kit.Observer, reviewRepository *ReviewRepository,
	config config.Config) *ReviewMiddleware {
	return &ReviewMiddleware{
		config:           config,
		observer:         observer,
		reviewRepository: reviewRepository,
	}
}

func (self *ReviewMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestCtx := ctx.Request().Context()
		requestProduct := product.RequestProduct(requestCtx)

		review, err := self.reviewRepository.GetByID(requestCtx, ctx.Param("review_id"))
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		if review == nil {
			return kit.HTTPErrInvalidRequest
		}

		if review.ProductID != requestProduct.ID {
			return kit.HTTPErrUnauthorized
		}

		ctx.SetRequest(ctx.Request().WithContext(context.WithValue(requestCtx, KeyRequestReview, review)))

		return next(ctx)
	}
}
