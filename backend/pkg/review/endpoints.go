package review

import (
	"net/http"
	"time"

	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
	"backend/pkg/product"
	"backend/pkg/util"

	kitUtil "github.com/neoxelox/kit/util"
)

const (
	REVIEW_ENDPOINTS_SEARCH_KEY       = "review:endpoints:search:"
	REVIEW_ENDPOINTS_SEARCH_TTL       = 10 * time.Minute
	REVIEW_ENDPOINTS_SEARCH_MAX_LIMIT = 100
)

type ReviewEndpoints struct {
	config           config.Config
	observer         *kit.Observer
	reviewRepository *ReviewRepository
	cache            *kit.Cache
}

func NewReviewEndpoints(observer *kit.Observer, reviewRepository *ReviewRepository,
	cache *kit.Cache, config config.Config) *ReviewEndpoints {
	return &ReviewEndpoints{
		config:           config,
		observer:         observer,
		reviewRepository: reviewRepository,
		cache:            cache,
	}
}

type ReviewEndpointsListReviewsRequest struct {
	Filters struct {
		Sources     []string   `query:"sources"`
		Releases    []string   `query:"releases"`
		Categories  []string   `query:"categories"`
		Keywords    []string   `query:"keywords"`
		Sentiments  []string   `query:"sentiments"`
		Emotions    []string   `query:"emotions"`
		Intentions  []string   `query:"intentions"`
		Languages   []string   `query:"languages"`
		SeenStartAt *time.Time `query:"seen_start_at"`
		SeenEndAt   *time.Time `query:"seen_end_at"`
	}
	Orders struct {
		Recency *string `query:"recency"`
	}
	Pagination struct {
		Limit *int    `query:"limit"`
		From  *string `query:"from"`
	}
}

type ReviewEndpointsListReviewsResponse struct {
	Reviews []ReviewPayload `json:"reviews"`
	Next    *string         `json:"next"`
}

func (self *ReviewEndpoints) ListReviews(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := ReviewEndpointsListReviewsRequest{}

	response := ReviewEndpointsListReviewsResponse{}
	err := self.cache.Get(requestCtx, REVIEW_ENDPOINTS_SEARCH_KEY+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.Filters.SeenStartAt != nil && request.Filters.SeenEndAt != nil &&
		request.Filters.SeenStartAt.After(*request.Filters.SeenEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	if request.Orders.Recency != nil {
		if *request.Orders.Recency != ReviewSearchOrdersDescending &&
			*request.Orders.Recency != ReviewSearchOrdersAscending {
			return kit.HTTPErrInvalidRequest
		}
	} else {
		request.Orders.Recency = kitUtil.Pointer(ReviewSearchOrdersDescending)
	}

	if request.Pagination.Limit != nil {
		if *request.Pagination.Limit > REVIEW_ENDPOINTS_SEARCH_MAX_LIMIT {
			return kit.HTTPErrInvalidRequest
		}
	} else {
		request.Pagination.Limit = kitUtil.Pointer(100)
	}

	page, err := self.reviewRepository.ListByProductID(requestCtx, requestProduct.ID, ReviewSearch{
		Filters: ReviewSearchFilters{
			Sources:     &request.Filters.Sources,
			Releases:    &request.Filters.Releases,
			Categories:  &request.Filters.Categories,
			Keywords:    &request.Filters.Keywords,
			Sentiments:  &request.Filters.Sentiments,
			Emotions:    &request.Filters.Emotions,
			Intentions:  &request.Filters.Intentions,
			Languages:   &request.Filters.Languages,
			SeenStartAt: request.Filters.SeenStartAt,
			SeenEndAt:   request.Filters.SeenEndAt,
		},
		Orders: ReviewSearchOrders{
			Recency: *request.Orders.Recency,
		},
		Pagination: util.Pagination[time.Time]{
			Limit: *request.Pagination.Limit,
			From:  util.CursorFromString[time.Time](request.Pagination.From),
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = ReviewEndpointsListReviewsResponse{}
	response.Reviews = make([]ReviewPayload, 0, len(page.Items))
	for _, review := range page.Items {
		response.Reviews = append(response.Reviews, *NewReviewPayload(review))
	}
	response.Next = nil
	next := util.CursorToString(page.Next)
	if next != nil {
		*next = url.QueryEscape(*next)
	}
	response.Next = next

	err = self.cache.Set(requestCtx, REVIEW_ENDPOINTS_SEARCH_KEY+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(REVIEW_ENDPOINTS_SEARCH_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type ReviewEndpointsGetReviewResponse struct {
	ReviewPayload
}

func (self *ReviewEndpoints) GetReview(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestReview := RequestReview(requestCtx)

	response := ReviewEndpointsGetReviewResponse{}
	response.ReviewPayload = *NewReviewPayload(*requestReview)

	return ctx.JSON(http.StatusOK, &response)
}

type ReviewEndpointsPutReviewQualityRequest struct {
	Quality int `json:"quality"`
}

type ReviewEndpointsPutReviewQualityResponse struct {
	Quality int `json:"quality"`
}

func (self *ReviewEndpoints) PutReviewQuality(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestReview := RequestReview(requestCtx)
	request := ReviewEndpointsPutReviewQualityRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.Quality < REVIEW_MIN_QUALITY || request.Quality > REVIEW_MAX_QUALITY {
		return kit.HTTPErrInvalidRequest
	}

	requestReview.Quality = &request.Quality

	err = self.reviewRepository.UpdateQuality(requestCtx, requestReview.ID, *requestReview.Quality)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := ReviewEndpointsPutReviewQualityResponse{}
	response.Quality = *requestReview.Quality

	return ctx.JSON(http.StatusOK, &response)
}
