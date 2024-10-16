package suggestion

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
	SUGGESTION_ENDPOINTS_SEARCH_KEY                = "suggestion:endpoints:search:"
	SUGGESTION_ENDPOINTS_SEARCH_TTL                = 10 * time.Minute
	SUGGESTION_ENDPOINTS_SEARCH_CONTENT_KEY        = "suggestion:endpoints:search:content:"
	SUGGESTION_ENDPOINTS_SEARCH_CONTENT_TTL        = 24 * time.Hour
	SUGGESTION_ENDPOINTS_SEARCH_MAX_LIMIT          = 100
	SUGGESTION_ENDPOINTS_SEARCH_MIN_CONTENT_LENGTH = 5
	SUGGESTION_ENDPOINTS_SEARCH_MAX_CONTENT_LENGTH = 250
)

type SuggestionEndpoints struct {
	config               config.Config
	observer             *kit.Observer
	suggestionRepository *SuggestionRepository
	userRepository       user.UserRepository
	engineService        *engine.EngineService
	cache                *kit.Cache
}

func NewSuggestionEndpoints(observer *kit.Observer, suggestionRepository *SuggestionRepository,
	userRepository user.UserRepository, engineService *engine.EngineService, cache *kit.Cache,
	config config.Config) *SuggestionEndpoints {
	return &SuggestionEndpoints{
		config:               config,
		observer:             observer,
		suggestionRepository: suggestionRepository,
		userRepository:       userRepository,
		engineService:        engineService,
		cache:                cache,
	}
}

type SuggestionEndpointsListSuggestionsRequest struct {
	Filters struct {
		Content          *string    `query:"content"`
		Sources          []string   `query:"sources"`
		Importances      []string   `query:"importances"`
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

type SuggestionEndpointsListSuggestionsResponse struct {
	Suggestions []SuggestionPayload `json:"suggestions"`
	Next        *string             `json:"next"`
}

func (self *SuggestionEndpoints) ListSuggestions(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := SuggestionEndpointsListSuggestionsRequest{}

	response := SuggestionEndpointsListSuggestionsResponse{}
	err := self.cache.Get(requestCtx, SUGGESTION_ENDPOINTS_SEARCH_KEY+requestProduct.ID+ctx.QueryString(), &response)
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

		if len(*request.Filters.Content) < SUGGESTION_ENDPOINTS_SEARCH_MIN_CONTENT_LENGTH ||
			len(*request.Filters.Content) > SUGGESTION_ENDPOINTS_SEARCH_MAX_CONTENT_LENGTH {
			return kit.HTTPErrInvalidRequest
		}
	}

	if request.Filters.Status != nil && !IsSuggestionSearchFiltersStatus(*request.Filters.Status) {
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
		if *request.Orders.Relevance != SuggestionSearchOrdersDescending &&
			*request.Orders.Relevance != SuggestionSearchOrdersAscending {
			return kit.HTTPErrInvalidRequest
		}
	} else {
		request.Orders.Relevance = kitUtil.Pointer(SuggestionSearchOrdersDescending)
	}

	if request.Pagination.Limit != nil {
		if *request.Pagination.Limit > SUGGESTION_ENDPOINTS_SEARCH_MAX_LIMIT {
			return kit.HTTPErrInvalidRequest
		}
	} else {
		request.Pagination.Limit = kitUtil.Pointer(100)
	}

	var embedding []float32
	if request.Filters.Content != nil {
		err := self.cache.Get(requestCtx, SUGGESTION_ENDPOINTS_SEARCH_CONTENT_KEY+*request.Filters.Content, &embedding)
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

				err = self.cache.Set(requestCtx, SUGGESTION_ENDPOINTS_SEARCH_CONTENT_KEY+*request.Filters.Content,
					embedding, kitUtil.Pointer(SUGGESTION_ENDPOINTS_SEARCH_CONTENT_TTL))
				if err != nil {
					return kit.HTTPErrServerGeneric.Cause(err)
				}
			} else {
				return kit.HTTPErrServerGeneric.Cause(err)
			}
		}
	}

	page, err := self.suggestionRepository.ListByProductID(requestCtx, requestProduct.ID, SuggestionSearch{
		Filters: SuggestionSearchFilters{
			Embedding:        &embedding,
			Sources:          &request.Filters.Sources,
			Importances:      &request.Filters.Importances,
			Releases:         &request.Filters.Releases,
			Categories:       &request.Filters.Categories,
			Assignees:        &request.Filters.Assignees,
			Status:           request.Filters.Status,
			FirstSeenStartAt: request.Filters.FirstSeenStartAt,
			FirstSeenEndAt:   request.Filters.FirstSeenEndAt,
			LastSeenStartAt:  request.Filters.LastSeenStartAt,
			LastSeenEndAt:    request.Filters.LastSeenEndAt,
		},
		Orders: SuggestionSearchOrders{
			Relevance: *request.Orders.Relevance,
		},
		Pagination: util.Pagination[SuggestionSearchCursor]{
			Limit: *request.Pagination.Limit,
			From:  util.CursorFromString[SuggestionSearchCursor](request.Pagination.From),
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = SuggestionEndpointsListSuggestionsResponse{}
	response.Suggestions = make([]SuggestionPayload, 0, len(page.Items))
	for _, suggestion := range page.Items {
		response.Suggestions = append(response.Suggestions, *NewSuggestionPayload(suggestion))
	}
	response.Next = nil
	next := util.CursorToString(page.Next)
	if next != nil {
		*next = url.QueryEscape(*next)
	}
	response.Next = next

	err = self.cache.Set(requestCtx, SUGGESTION_ENDPOINTS_SEARCH_KEY+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(SUGGESTION_ENDPOINTS_SEARCH_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type SuggestionEndpointsGetSuggestionResponse struct {
	SuggestionPayload
}

func (self *SuggestionEndpoints) GetSuggestion(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestSuggestion := RequestSuggestion(requestCtx)

	response := SuggestionEndpointsGetSuggestionResponse{}
	response.SuggestionPayload = *NewSuggestionPayload(*requestSuggestion)

	return ctx.JSON(http.StatusOK, &response)
}

type SuggestionEndpointsListSuggestionFeedbacksRequest struct {
	From *string `query:"from"`
}

type SuggestionEndpointsListSuggestionFeedbacksResponse struct {
	Feedbacks []feedback.FeedbackPayload `json:"feedbacks"`
	Next      *string                    `json:"next"`
}

func (self *SuggestionEndpoints) ListSuggestionFeedbacks(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestSuggestion := RequestSuggestion(requestCtx)
	request := SuggestionEndpointsListSuggestionFeedbacksRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	page, err := self.suggestionRepository.ListFeedbacks(requestCtx, requestSuggestion.ID, util.Pagination[time.Time]{
		Limit: 100,
		From:  util.CursorFromString[time.Time](request.From),
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := SuggestionEndpointsListSuggestionFeedbacksResponse{}
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

type SuggestionEndpointsPutSuggestionAssigneeRequest struct {
	AssigneeID *string `json:"assignee_id"`
}

type SuggestionEndpointsPutSuggestionAssigneeResponse struct {
	AssigneeID *string `json:"assignee_id"`
}

func (self *SuggestionEndpoints) PutSuggestionAssignee(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestOrganization := organization.RequestOrganization(requestCtx)
	requestProduct := product.RequestProduct(requestCtx)
	requestSuggestion := RequestSuggestion(requestCtx)
	request := SuggestionEndpointsPutSuggestionAssigneeRequest{}

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

	requestSuggestion.AssigneeID = request.AssigneeID

	err = self.suggestionRepository.UpdateAssignee(requestCtx, requestSuggestion.ID, requestSuggestion.AssigneeID)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	keys, err := self.cache.Find(requestCtx, SUGGESTION_ENDPOINTS_SEARCH_KEY+requestProduct.ID+"*")
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	for _, key := range keys {
		err = self.cache.Delete(requestCtx, key)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}
	}

	response := SuggestionEndpointsPutSuggestionAssigneeResponse{}
	response.AssigneeID = requestSuggestion.AssigneeID

	return ctx.JSON(http.StatusOK, &response)
}

type SuggestionEndpointsPutSuggestionQualityRequest struct {
	Quality int `json:"quality"`
}

type SuggestionEndpointsPutSuggestionQualityResponse struct {
	Quality int `json:"quality"`
}

func (self *SuggestionEndpoints) PutSuggestionQuality(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestSuggestion := RequestSuggestion(requestCtx)
	request := SuggestionEndpointsPutSuggestionQualityRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.Quality < SUGGESTION_MIN_QUALITY || request.Quality > SUGGESTION_MAX_QUALITY {
		return kit.HTTPErrInvalidRequest
	}

	requestSuggestion.Quality = &request.Quality

	err = self.suggestionRepository.UpdateQuality(requestCtx, requestSuggestion.ID, *requestSuggestion.Quality)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := SuggestionEndpointsPutSuggestionQualityResponse{}
	response.Quality = *requestSuggestion.Quality

	return ctx.JSON(http.StatusOK, &response)
}

type SuggestionEndpointsPutSuggestionArchivedRequest struct {
	Archived bool `json:"archived"`
}

type SuggestionEndpointsPutSuggestionArchivedResponse struct {
	Archived bool `json:"archived"`
}

func (self *SuggestionEndpoints) PutSuggestionArchived(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	requestSuggestion := RequestSuggestion(requestCtx)
	request := SuggestionEndpointsPutSuggestionArchivedRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.Archived {
		requestSuggestion.ArchivedAt = kitUtil.Pointer(time.Now())
	} else {
		requestSuggestion.ArchivedAt = nil
	}

	err = self.suggestionRepository.UpdateArchivedAt(requestCtx, requestSuggestion.ID, requestSuggestion.ArchivedAt)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	keys, err := self.cache.Find(requestCtx, SUGGESTION_ENDPOINTS_SEARCH_KEY+requestProduct.ID+"*")
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	for _, key := range keys {
		err = self.cache.Delete(requestCtx, key)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}
	}

	response := SuggestionEndpointsPutSuggestionArchivedResponse{}
	response.Archived = request.Archived

	return ctx.JSON(http.StatusOK, &response)
}
