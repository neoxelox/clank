package issue

import (
	"context"

	"backend/pkg/config"
	"backend/pkg/product"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"
)

var (
	KeyRequestIssue kit.Key = kit.KeyBase + "request:issue"
)

func RequestIssue(ctx context.Context) *Issue {
	return ctx.Value(KeyRequestIssue).(*Issue) // nolint:forcetypeassert,errcheck
}

type IssueMiddleware struct {
	config          config.Config
	observer        *kit.Observer
	issueRepository *IssueRepository
}

func NewIssueMiddleware(observer *kit.Observer, issueRepository *IssueRepository,
	config config.Config) *IssueMiddleware {
	return &IssueMiddleware{
		config:          config,
		observer:        observer,
		issueRepository: issueRepository,
	}
}

func (self *IssueMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestCtx := ctx.Request().Context()
		requestProduct := product.RequestProduct(requestCtx)

		issue, err := self.issueRepository.GetByID(requestCtx, ctx.Param("issue_id"))
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		if issue == nil {
			return kit.HTTPErrInvalidRequest
		}

		if issue.ProductID != requestProduct.ID {
			return kit.HTTPErrUnauthorized
		}

		ctx.SetRequest(ctx.Request().WithContext(context.WithValue(requestCtx, KeyRequestIssue, issue)))

		return next(ctx)
	}
}
