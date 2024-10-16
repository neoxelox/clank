package metric

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"
	kitUtil "github.com/neoxelox/kit/util"

	"backend/pkg/config"
	"backend/pkg/product"
)

const (
	METRIC_ENDPOINTS_KEY = "metric:endpoints:"
	METRIC_ENDPOINTS_TTL = 10 * time.Minute
)

type MetricEndpoints struct {
	config           config.Config
	observer         *kit.Observer
	metricRepository *MetricRepository
	cache            *kit.Cache
}

func NewMetricEndpoints(observer *kit.Observer, metricRepository *MetricRepository,
	cache *kit.Cache, config config.Config) *MetricEndpoints {
	return &MetricEndpoints{
		config:           config,
		observer:         observer,
		metricRepository: metricRepository,
		cache:            cache,
	}
}

type MetricEndpointsGetRequest struct {
	PeriodStartAt *time.Time `query:"period_start_at"`
	PeriodEndAt   *time.Time `query:"period_end_at"`
}

type MetricEndpointsGetResponse struct {
}

type MetricEndpointsGetIssueCountRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetIssueCountResponse struct {
	MetricEndpointsGetResponse
	ActiveIssues   int `json:"active_issues"`
	ArchivedIssues int `json:"archived_issues"`
	NewIssues      int `json:"new_issues"`
	Feedbacks      int `json:"feedbacks"`
}

func (self *MetricEndpoints) GetIssueCount(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetIssueCountRequest{}

	response := MetricEndpointsGetIssueCountResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"issue:count:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetIssueCount(requestCtx, IssueCountParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetIssueCountResponse{}
	response.ActiveIssues = metric.ActiveIssues
	response.ArchivedIssues = metric.ArchivedIssues
	response.NewIssues = metric.NewIssues
	response.Feedbacks = metric.Feedbacks

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"issue:count:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetIssueSourcesRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetIssueSourcesResponse struct {
	MetricEndpointsGetResponse
	Sources map[string]int `json:"sources"`
}

func (self *MetricEndpoints) GetIssueSources(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetIssueSourcesRequest{}

	response := MetricEndpointsGetIssueSourcesResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"issue:sources:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetIssueSources(requestCtx, IssueSourcesParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetIssueSourcesResponse{}
	response.Sources = metric.Sources

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"issue:sources:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetIssueSeveritiesRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetIssueSeveritiesResponse struct {
	MetricEndpointsGetResponse
	Severities map[string]int `json:"severities"`
}

func (self *MetricEndpoints) GetIssueSeverities(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetIssueSeveritiesRequest{}

	response := MetricEndpointsGetIssueSeveritiesResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"issue:severities:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetIssueSeverities(requestCtx, IssueSeveritiesParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetIssueSeveritiesResponse{}
	response.Severities = metric.Severities

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"issue:severities:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetIssueCategoriesRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetIssueCategoriesResponse struct {
	MetricEndpointsGetResponse
	Categories map[string]int `json:"categories"`
}

func (self *MetricEndpoints) GetIssueCategories(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetIssueCategoriesRequest{}

	response := MetricEndpointsGetIssueCategoriesResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"issue:categories:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetIssueCategories(requestCtx, IssueCategoriesParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetIssueCategoriesResponse{}
	response.Categories = metric.Categories

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"issue:categories:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetIssueReleasesRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetIssueReleasesResponse struct {
	MetricEndpointsGetResponse
	Releases map[string]int `json:"releases"`
}

func (self *MetricEndpoints) GetIssueReleases(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetIssueReleasesRequest{}

	response := MetricEndpointsGetIssueReleasesResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"issue:releases:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetIssueReleases(requestCtx, IssueReleasesParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetIssueReleasesResponse{}
	response.Releases = metric.Releases

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"issue:releases:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetSuggestionCountRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetSuggestionCountResponse struct {
	MetricEndpointsGetResponse
	ActiveSuggestions   int `json:"active_suggestions"`
	ArchivedSuggestions int `json:"archived_suggestions"`
	NewSuggestions      int `json:"new_suggestions"`
	Feedbacks           int `json:"feedbacks"`
}

func (self *MetricEndpoints) GetSuggestionCount(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetSuggestionCountRequest{}

	response := MetricEndpointsGetSuggestionCountResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"suggestion:count:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetSuggestionCount(requestCtx, SuggestionCountParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetSuggestionCountResponse{}
	response.ActiveSuggestions = metric.ActiveSuggestions
	response.ArchivedSuggestions = metric.ArchivedSuggestions
	response.NewSuggestions = metric.NewSuggestions
	response.Feedbacks = metric.Feedbacks

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"suggestion:count:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetSuggestionSourcesRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetSuggestionSourcesResponse struct {
	MetricEndpointsGetResponse
	Sources map[string]int `json:"sources"`
}

func (self *MetricEndpoints) GetSuggestionSources(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetSuggestionSourcesRequest{}

	response := MetricEndpointsGetSuggestionSourcesResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"suggestion:sources:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetSuggestionSources(requestCtx, SuggestionSourcesParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetSuggestionSourcesResponse{}
	response.Sources = metric.Sources

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"suggestion:sources:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetSuggestionImportancesRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetSuggestionImportancesResponse struct {
	MetricEndpointsGetResponse
	Importances map[string]int `json:"importances"`
}

func (self *MetricEndpoints) GetSuggestionImportances(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetSuggestionImportancesRequest{}

	response := MetricEndpointsGetSuggestionImportancesResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"suggestion:importances:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetSuggestionImportances(requestCtx, SuggestionImportancesParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetSuggestionImportancesResponse{}
	response.Importances = metric.Importances

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"suggestion:importances:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetSuggestionCategoriesRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetSuggestionCategoriesResponse struct {
	MetricEndpointsGetResponse
	Categories map[string]int `json:"categories"`
}

func (self *MetricEndpoints) GetSuggestionCategories(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetSuggestionCategoriesRequest{}

	response := MetricEndpointsGetSuggestionCategoriesResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"suggestion:categories:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetSuggestionCategories(requestCtx, SuggestionCategoriesParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetSuggestionCategoriesResponse{}
	response.Categories = metric.Categories

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"suggestion:categories:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetSuggestionReleasesRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetSuggestionReleasesResponse struct {
	MetricEndpointsGetResponse
	Releases map[string]int `json:"releases"`
}

func (self *MetricEndpoints) GetSuggestionReleases(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetSuggestionReleasesRequest{}

	response := MetricEndpointsGetSuggestionReleasesResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"suggestion:releases:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetSuggestionReleases(requestCtx, SuggestionReleasesParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetSuggestionReleasesResponse{}
	response.Releases = metric.Releases

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"suggestion:releases:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetReviewSentimentsRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetReviewSentimentsResponse struct {
	MetricEndpointsGetResponse
	Sentiments map[string]int `json:"sentiments"`
}

func (self *MetricEndpoints) GetReviewSentiments(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetReviewSentimentsRequest{}

	response := MetricEndpointsGetReviewSentimentsResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"review:sentiments:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetReviewSentiments(requestCtx, ReviewSentimentsParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetReviewSentimentsResponse{}
	response.Sentiments = metric.Sentiments

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"review:sentiments:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetReviewSourcesRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetReviewSourcesResponse struct {
	MetricEndpointsGetResponse
	Sources map[string]int `json:"sources"`
}

func (self *MetricEndpoints) GetReviewSources(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetReviewSourcesRequest{}

	response := MetricEndpointsGetReviewSourcesResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"review:sources:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetReviewSources(requestCtx, ReviewSourcesParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetReviewSourcesResponse{}
	response.Sources = metric.Sources

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"review:sources:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetReviewIntentionsRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetReviewIntentionsResponse struct {
	MetricEndpointsGetResponse
	Intentions map[string]int `json:"intentions"`
}

func (self *MetricEndpoints) GetReviewIntentions(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetReviewIntentionsRequest{}

	response := MetricEndpointsGetReviewIntentionsResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"review:intentions:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetReviewIntentions(requestCtx, ReviewIntentionsParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetReviewIntentionsResponse{}
	response.Intentions = metric.Intentions

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"review:intentions:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetReviewEmotionsRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetReviewEmotionsResponse struct {
	MetricEndpointsGetResponse
	Emotions map[string]int `json:"emotions"`
}

func (self *MetricEndpoints) GetReviewEmotions(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetReviewEmotionsRequest{}

	response := MetricEndpointsGetReviewEmotionsResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"review:emotions:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetReviewEmotions(requestCtx, ReviewEmotionsParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetReviewEmotionsResponse{}
	response.Emotions = metric.Emotions

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"review:emotions:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetReviewCategoriesRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetReviewCategoriesResponse struct {
	MetricEndpointsGetResponse
	Categories map[string]int `json:"categories"`
}

func (self *MetricEndpoints) GetReviewCategories(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetReviewCategoriesRequest{}

	response := MetricEndpointsGetReviewCategoriesResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"review:categories:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetReviewCategories(requestCtx, ReviewCategoriesParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetReviewCategoriesResponse{}
	response.Categories = metric.Categories

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"review:categories:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetReviewReleasesRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetReviewReleasesResponse struct {
	MetricEndpointsGetResponse
	Releases map[string]int `json:"releases"`
}

func (self *MetricEndpoints) GetReviewReleases(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetReviewReleasesRequest{}

	response := MetricEndpointsGetReviewReleasesResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"review:releases:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetReviewReleases(requestCtx, ReviewReleasesParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetReviewReleasesResponse{}
	response.Releases = metric.Releases

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"review:releases:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetReviewKeywordsRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetReviewKeywordsResponse struct {
	MetricEndpointsGetResponse
	Positive map[string]int `json:"positive"`
	Neutral  map[string]int `json:"neutral"`
	Negative map[string]int `json:"negative"`
}

func (self *MetricEndpoints) GetReviewKeywords(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetReviewKeywordsRequest{}

	response := MetricEndpointsGetReviewKeywordsResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"review:keywords:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetReviewKeywords(requestCtx, ReviewKeywordsParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetReviewKeywordsResponse{}
	response.Positive = metric.Positive
	response.Neutral = metric.Neutral
	response.Negative = metric.Negative

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"review:keywords:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetNetPromoterScoreRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetNetPromoterScoreResponse struct {
	MetricEndpointsGetResponse
	Score float64 `json:"score"`
}

func (self *MetricEndpoints) GetNetPromoterScore(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetNetPromoterScoreRequest{}

	response := MetricEndpointsGetNetPromoterScoreResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"nps:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetNetPromoterScore(requestCtx, NetPromoterScoreParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetNetPromoterScoreResponse{}
	response.Score = metric.Score

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"nps:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}

type MetricEndpointsGetCustomerSatisfactionScoreRequest struct {
	MetricEndpointsGetRequest
}

type MetricEndpointsGetCustomerSatisfactionScoreResponse struct {
	MetricEndpointsGetResponse
	Score float64 `json:"score"`
}

func (self *MetricEndpoints) GetCustomerSatisfactionScore(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)
	request := MetricEndpointsGetCustomerSatisfactionScoreRequest{}

	response := MetricEndpointsGetCustomerSatisfactionScoreResponse{}
	err := self.cache.Get(requestCtx,
		METRIC_ENDPOINTS_KEY+"csat:"+requestProduct.ID+ctx.QueryString(), &response)
	if err != nil && !kit.ErrCacheMiss.Is(err) {
		return kit.HTTPErrServerGeneric.Cause(err)
	} else if err == nil {
		return ctx.JSON(http.StatusOK, &response)
	}

	err = ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.PeriodStartAt != nil && request.PeriodEndAt != nil &&
		request.PeriodStartAt.After(*request.PeriodEndAt) {
		return kit.HTTPErrInvalidRequest
	}

	metric, err := self.metricRepository.GetCustomerSatisfactionScore(requestCtx, CustomerSatisfactionScoreParams{
		Params: Params{
			ProductID:     requestProduct.ID,
			PeriodStartAt: request.PeriodStartAt,
			PeriodEndAt:   request.PeriodEndAt,
		},
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response = MetricEndpointsGetCustomerSatisfactionScoreResponse{}
	response.Score = metric.Score

	err = self.cache.Set(requestCtx,
		METRIC_ENDPOINTS_KEY+"csat:"+requestProduct.ID+ctx.QueryString(),
		&response, kitUtil.Pointer(METRIC_ENDPOINTS_TTL))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, &response)
}
