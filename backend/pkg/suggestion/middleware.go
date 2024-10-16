package suggestion

import (
	"context"

	"backend/pkg/config"
	"backend/pkg/product"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"
)

var (
	KeyRequestSuggestion kit.Key = kit.KeyBase + "request:suggestion"
)

func RequestSuggestion(ctx context.Context) *Suggestion {
	return ctx.Value(KeyRequestSuggestion).(*Suggestion) // nolint:forcetypeassert,errcheck
}

type SuggestionMiddleware struct {
	config               config.Config
	observer             *kit.Observer
	suggestionRepository *SuggestionRepository
}

func NewSuggestionMiddleware(observer *kit.Observer, suggestionRepository *SuggestionRepository,
	config config.Config) *SuggestionMiddleware {
	return &SuggestionMiddleware{
		config:               config,
		observer:             observer,
		suggestionRepository: suggestionRepository,
	}
}

func (self *SuggestionMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestCtx := ctx.Request().Context()
		requestProduct := product.RequestProduct(requestCtx)

		suggestion, err := self.suggestionRepository.GetByID(requestCtx, ctx.Param("suggestion_id"))
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		if suggestion == nil {
			return kit.HTTPErrInvalidRequest
		}

		if suggestion.ProductID != requestProduct.ID {
			return kit.HTTPErrUnauthorized
		}

		ctx.SetRequest(ctx.Request().WithContext(context.WithValue(requestCtx, KeyRequestSuggestion, suggestion)))

		return next(ctx)
	}
}
