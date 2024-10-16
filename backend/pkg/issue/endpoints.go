package issue

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"
	kitUtil "github.com/neoxelox/kit/util"

	"backend/pkg/config"
	"backend/pkg/engine"
	"backend/pkg/feedback"
	"backend/pkg/organization"
	"backend/pkg/product"
	"backend/pkg/user"
	"backend/pkg/util"
)

const (
	ISSUE_ENDPOINTS_SEARCH_KEY                = "issue:endpoints:search:"
	ISSUE_ENDPOINTS_SEARCH_TTL                = 10 * time.Minute
	ISSUE_ENDPOINTS_SEARCH_CONTENT_KEY        = "issue:endpoints:search:content:"
	ISSUE_ENDPOINTS_SEARCH_CONTENT_TTL        = 24 * time.Hour
	ISSUE_ENDPOINTS_SEARCH_MAX_LIMIT          = 100
	ISSUE_ENDPOINTS_SEARCH_MIN_CONTENT_LENGTH = 5
	ISSUE_ENDPOINTS_SEARCH_MAX_CONTENT_LENGTH = 250
)

type IssueEndpoints struct {
	config          config.Config
	observer        *kit.Observer
	issueRepository *IssueRepository
	userRepository  user.UserRepository
	engineService   *engine.EngineService
	cache           *kit.Cache
}

func NewIssueEndpoints(observer *kit.Observer, issueRepository *IssueRepository, userRepository user.UserRepository,
	engineService *engine.EngineService, cache *kit.Cache, config config.Config) *IssueEndpoints {
	return &IssueEndpoints{
		config:          config,
		observer:        observer,
		issueRepository: issueRepository,
		userRepository:  userRepository,
		engineService:   engineService,
		cache:           cache,
	}
}

type IssueEndpointsListIssuesRequest struct {
	Filters struct {
		Content          *string    `query:"content"`
		Sources          []string   `query:"sources"`
		Severities       []string   `query:"severities"`
		Releases         []string   `query:"releases"`
		Categories       []string   `query:"categories"`
		Assignees        []string   `query:"assignees"`
		Status           *string    `query:"status"`
		FirstSeenStartAt *time.Time `query:"first_seen_start_at"`
		FirstSeenEndAt   *time.Time `query:"first_seen_end_at"`
		LastSeenStartAt  *time.Time `query:"last_seen_start_at"`
		LastSeenEndAt    *time.Time `query:"last_seen_end_at"`
	}
	Orders struct {
		Relevance *string `query:"relevance"`
	}
	Pagination struct {
		Limit *int    `query:"limit"`
		From  *string `query:"from"`
	}
}

type IssueEndpointsListIssuesResponse struct {
	Issues []IssuePayload `json:"issues"`
	Next   *string        `json:"next"`
}

func (self *IssueEndpoints) ListIssues(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := IssueEndpointsListIssuesRequest{}

	response := IssueEndpointsListIssuesResponse{}
	err := self.cache.Get(requestCtx, ISSUE_ENDPOINTS_SEARCH_KEY+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.Filters.Content != nil {
		*request.Filters.Content = strings.Trim(*request.Filters.Content, `"`)

		if len(*request.Filters.Content) < ISSUE_ENDPOINTS_SEARCH_MIN_CONTENT_LENGTH ||
			len(*request.Filters.Content) > ISSUE_ENDPOINTS_SEARCH_MAX_CONTENT_LENGTH {
			return kit.HTTPErrInvalidRequest
		}
	}

	if request.Filters.Status != nil && !IsIssueSearchFiltersStatus(*request.Filters.Status) {
		return kit.HTTPErrInvalidRequest
	}

	if request.Filters.FirstSeenStartAt != nil && request.Filters.FirstSeenEndAt != nil &&
		request.Filters.FirstSeenStartAt.After(*request.Filters.FirstSeenEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	if request.Filters.LastSeenStartAt != nil && request.Filters.LastSeenEndAt != nil &&
		request.Filters.LastSeenStartAt.After(*request.Filters.LastSeenEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	if request.Orders.Relevance != nil {
		if *request.Orders.Relevance != IssueSearchOrdersDescending &&
			*request.Orders.Relevance != IssueSearchOrdersAscending {
			return kit.HTTPErrInvalidRequest
		}
	} else {
		request.Orders.Relevance = kitUtil.Pointer(IssueSearchOrdersDescending)
	}

	if request.Pagination.Limit != nil {
		if *request.Pagination.Limit > ISSUE_ENDPOINTS_SEARCH_MAX_LIMIT {
			return kit.HTTPErrInvalidRequest
		}
	} else {
		request.Pagination.Limit = kitUtil.Pointer(100)
	}

	var embedding []float32
	if request.Filters.Content != nil {
		err := self.cache.Get(requestCtx, ISSUE_ENDPOINTS_SEARCH_CONTENT_KEY+*request.Filters.Content, &embedding)
		if err != nil {
			if kit.ErrCacheMiss.Is(err) {
				result, err := self.engineService.ComputeEmbedding(requestCtx,
					engine.EngineServiceComputeEmbeddingParams{
						Text: *request.Filters.Content,
					})
				if err != nil {
					return kit.HTTPErrServerGeneric.Cause(err)
				}

				embedding = result.Embedding

				err = self.cache.Set(requestCtx, ISSUE_ENDPOINTS_SEARCH_CONTENT_KEY+*request.Filters.Content,
					embedding, kitUtil.Pointer(ISSUE_ENDPOINTS_SEARCH_CONTENT_TTL))
				if err != nil {
					return kit.HTTPErrServerGeneric.Cause(err)
				}
			} else {
				return kit.HTTPErrServerGeneric.Cause(err)
			}
		}
	}

	page, err := self.issueRepository.ListByProductID(requestCtx, requestProduct.ID, IssueSearch{
		Filters: IssueSearchFilters{
			Embedding:        &embedding,
			Sources:          &request.Filters.Sources,
			Severities:       &request.Filters.Severities,
			Releases:         &request.Filters.Releases,
			Categories:       &request.Filters.Categories,
			Assignees:        &request.Filters.Assignees,
			Status:           request.Filters.Status,
			FirstSeenStartAt: request.Filters.FirstSeenStartAt,
			FirstSeenEndAt:   request.Filters.FirstSeenEndAt,
			LastSeenStartAt:  request.Filters.LastSeenStartAt,
			LastSeenEndAt:    request.Filters.LastSeenEndAt,
		},
		Orders: IssueSearchOrders{
			Relevance: *request.Orders.Relevance,
		},
		Pagination: util.Pagination[IssueSearchCursor]{
			Limit: *request.Pagination.Limit,
			From:  util.CursorFromString[IssueSearchCursor](request.Pagination.From),
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = IssueEndpointsListIssuesResponse{}
	response.Issues = make([]IssuePayload, 0, len(page.Items))
	for _, issue := range page.Items {
		response.Issues = append(response.Issues, *NewIssuePayload(issue))
	}
	response.Next = nil
	next := util.CursorToString(page.Next)
	if next != nil {
		*next = url.QueryEscape(*next)
	}
	response.Next = next

	err = self.cache.Set(requestCtx, ISSUE_ENDPOINTS_SEARCH_KEY+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(ISSUE_ENDPOINTS_SEARCH_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type IssueEndpointsGetIssueResponse struct {
	IssuePayload
}

func (self *IssueEndpoints) GetIssue(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestIssue := RequestIssue(requestCtx)

	response := IssueEndpointsGetIssueResponse{}
	response.IssuePayload = *NewIssuePayload(*requestIssue)

	return ctx.JSON(http.StatusOK, &response)
}

type IssueEndpointsListIssueFeedbacksRequest struct {
	From *string `query:"from"`
}

type IssueEndpointsListIssueFeedbacksResponse struct {
	Feedbacks []feedback.FeedbackPayload `json:"feedbacks"`
	Next      *string                    `json:"next"`
}

func (self *IssueEndpoints) ListIssueFeedbacks(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestIssue := RequestIssue(requestCtx)
	request := IssueEndpointsListIssueFeedbacksRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	page, err := self.issueRepository.ListFeedbacks(requestCtx, requestIssue.ID, util.Pagination[time.Time]{
		Limit: 100,
		From:  util.CursorFromString[time.Time](request.From),
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := IssueEndpointsListIssueFeedbacksResponse{}
	response.Feedbacks = make([]feedback.FeedbackPayload, 0, len(page.Items))
	for _, _feedback := range page.Items {
		response.Feedbacks = append(response.Feedbacks, *feedback.NewFeedbackPayload(_feedback))
	}
	response.Next = nil
	next := util.CursorToString(page.Next)
	if next != nil {
		*next = url.QueryEscape(*next)
	}
	response.Next = next

	return ctx.JSON(http.StatusOK, &response)
}

type IssueEndpointsPutIssueAssigneeRequest struct {
	AssigneeID *string `json:"assignee_id"`
}

type IssueEndpointsPutIssueAssigneeResponse struct {
	AssigneeID *string `json:"assignee_id"`
}

func (self *IssueEndpoints) PutIssueAssignee(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestOrganization := organization.RequestOrganization(requestCtx)
	requestProduct := product.RequestProduct(requestCtx)
	requestIssue := RequestIssue(requestCtx)
	request := IssueEndpointsPutIssueAssigneeRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.AssigneeID != nil {
		user, err := self.userRepository.GetByID(requestCtx, *request.AssigneeID)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		if user == nil {
			return kit.HTTPErrInvalidRequest
		}

		if user.OrganizationID != requestOrganization.ID {
			return kit.HTTPErrUnauthorized
		}

		if user.DeletedAt != nil {
			return kit.HTTPErrUnauthorized
		}
	}

	requestIssue.AssigneeID = request.AssigneeID

	err = self.issueRepository.UpdateAssignee(requestCtx, requestIssue.ID, requestIssue.AssigneeID)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	keys, err := self.cache.Find(requestCtx, ISSUE_ENDPOINTS_SEARCH_KEY+requestProduct.ID+"*")
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	for _, key := range keys {
		err = self.cache.Delete(requestCtx, key)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}
	}

	response := IssueEndpointsPutIssueAssigneeResponse{}
	response.AssigneeID = requestIssue.AssigneeID

	return ctx.JSON(http.StatusOK, &response)
}

type IssueEndpointsPutIssueQualityRequest struct {
	Quality int `json:"quality"`
}

type IssueEndpointsPutIssueQualityResponse struct {
	Quality int `json:"quality"`
}

func (self *IssueEndpoints) PutIssueQuality(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestIssue := RequestIssue(requestCtx)
	request := IssueEndpointsPutIssueQualityRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.Quality < ISSUE_MIN_QUALITY || request.Quality > ISSUE_MAX_QUALITY {
		return kit.HTTPErrInvalidRequest
	}

	requestIssue.Quality = &request.Quality

	err = self.issueRepository.UpdateQuality(requestCtx, requestIssue.ID, *requestIssue.Quality)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := IssueEndpointsPutIssueQualityResponse{}
	response.Quality = *requestIssue.Quality

	return ctx.JSON(http.StatusOK, &response)
}

type IssueEndpointsPutIssueArchivedRequest struct {
	Archived bool `json:"archived"`
}

type IssueEndpointsPutIssueArchivedResponse struct {
	Archived bool `json:"archived"`
}

func (self *IssueEndpoints) PutIssueArchived(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	requestIssue := RequestIssue(requestCtx)
	request := IssueEndpointsPutIssueArchivedRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.Archived {
		requestIssue.ArchivedAt = kitUtil.Pointer(time.Now())
	} else {
		requestIssue.ArchivedAt = nil
	}

	err = self.issueRepository.UpdateArchivedAt(requestCtx, requestIssue.ID, requestIssue.ArchivedAt)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	keys, err := self.cache.Find(requestCtx, ISSUE_ENDPOINTS_SEARCH_KEY+requestProduct.ID+"*")
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	for _, key := range keys {
		err = self.cache.Delete(requestCtx, key)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}
	}

	response := IssueEndpointsPutIssueArchivedResponse{}
	response.Archived = request.Archived

	return ctx.JSON(http.StatusOK, &response)
}
