package collector

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"
	"github.com/neoxelox/kit/util"
	"github.com/rs/xid"

	"backend/pkg/config"
	"backend/pkg/organization"
	"backend/pkg/product"
)

type CollectorEndpoints struct {
	config              config.Config
	observer            *kit.Observer
	collectorRepository *CollectorRepository
	enqueuer            *kit.Enqueuer
}

func NewCollectorEndpoints(observer *kit.Observer, collectorRepository *CollectorRepository,
	enqueuer *kit.Enqueuer, config config.Config) *CollectorEndpoints {
	return &CollectorEndpoints{
		config:              config,
		observer:            observer,
		collectorRepository: collectorRepository,
		enqueuer:            enqueuer,
	}
}

type CollectorEndpointsListCollectorsResponse struct {
	Collectors []CollectorPayload `json:"collectors"`
}

func (self *CollectorEndpoints) ListCollectors(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)

	collectors, err := self.collectorRepository.ListByProductID(requestCtx, requestProduct.ID)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := CollectorEndpointsListCollectorsResponse{}
	response.Collectors = make([]CollectorPayload, 0, len(collectors))
	for _, collector := range collectors {
		if collector.DeletedAt != nil {
			continue
		}

		response.Collectors = append(response.Collectors, *NewCollectorPayload(collector))
	}

	return ctx.JSON(http.StatusOK, &response)
}

type CollectorEndpointsPostTrustpilotCollectorRequest struct {
	Domain string `json:"domain"`
}

type CollectorEndpointsPostPlayStoreCollectorRequest struct {
	AppID string `json:"app_id"`
}

type CollectorEndpointsPostAppStoreCollectorRequest struct {
	AppID string `json:"app_id"`
}

type CollectorEndpointsPostAmazonCollectorRequest struct {
	ASIN string `json:"asin"`
}

type CollectorEndpointsPostIAgoraCollectorRequest struct {
	Institution string `json:"institution"`
}

type CollectorEndpointsPostWebhookCollectorRequest struct {
}

type CollectorEndpointsPostWidgetCollectorRequest struct {
}

type CollectorEndpointsPostCollectorRequest struct {
	Type string `json:"type"`
}

type CollectorEndpointsPostCollectorResponse struct {
	CollectorPayload
}

func (self *CollectorEndpoints) PostCollector(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestOrganization := organization.RequestOrganization(requestCtx)
	requestProduct := product.RequestProduct(requestCtx)
	var requestRaw json.RawMessage
	request := CollectorEndpointsPostCollectorRequest{}

	err := ctx.Bind(&requestRaw)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	err = json.Unmarshal(requestRaw, &request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if !IsCollectorType(request.Type) {
		return kit.HTTPErrInvalidRequest
	}

	// TODO: Create an organization middleware that limits the usage of features by plan
	if requestOrganization.Plan == organization.OrganizationPlanTrial ||
		requestOrganization.Plan == organization.OrganizationPlanDemo {
		exists, err := self.collectorRepository.ExistsByOrganizationID(requestCtx, requestOrganization.ID)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		if exists {
			return kit.HTTPErrInvalidRequest
		}
	}

	collectors, err := self.collectorRepository.ListByProductIDAndType(requestCtx, requestProduct.ID, request.Type)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	var collector *Collector
	var settings any
	var jobdata any
	switch request.Type {
	case CollectorTypeTrustpilot:
		var _request CollectorEndpointsPostTrustpilotCollectorRequest
		err := json.Unmarshal(requestRaw, &_request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		if len(_request.Domain) == 0 {
			return kit.HTTPErrInvalidRequest
		}

		settings = TrustpilotCollectorSettings{
			Domain: _request.Domain,
		}
		jobdata = TrustpilotCollectorJobdata{
			LastDispatchedAt:    nil,
			LastDispatchedTasks: []string{},
			Cost:                0,
		}

		for _, _collector := range collectors {
			if _collector.Settings.(TrustpilotCollectorSettings).Domain == _request.Domain {
				collector = util.Pointer(_collector)
				break
			}
		}

	case CollectorTypePlayStore:
		var _request CollectorEndpointsPostPlayStoreCollectorRequest
		err := json.Unmarshal(requestRaw, &_request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		if len(_request.AppID) == 0 {
			return kit.HTTPErrInvalidRequest
		}

		settings = PlayStoreCollectorSettings{
			AppID: _request.AppID,
		}
		jobdata = PlayStoreCollectorJobdata{
			LastDispatchedAt:    nil,
			LastDispatchedTasks: []string{},
			Cost:                0,
		}

		for _, _collector := range collectors {
			if _collector.Settings.(PlayStoreCollectorSettings).AppID == _request.AppID {
				collector = util.Pointer(_collector)
				break
			}
		}

	case CollectorTypeAppStore:
		var _request CollectorEndpointsPostAppStoreCollectorRequest
		err := json.Unmarshal(requestRaw, &_request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		if len(_request.AppID) == 0 {
			return kit.HTTPErrInvalidRequest
		}

		settings = AppStoreCollectorSettings{
			AppID: _request.AppID,
		}
		jobdata = AppStoreCollectorJobdata{
			LastDispatchedAt:    nil,
			LastDispatchedTasks: []string{},
			Cost:                0,
		}

		for _, _collector := range collectors {
			if _collector.Settings.(AppStoreCollectorSettings).AppID == _request.AppID {
				collector = util.Pointer(_collector)
				break
			}
		}

	case CollectorTypeAmazon:
		var _request CollectorEndpointsPostAmazonCollectorRequest
		err := json.Unmarshal(requestRaw, &_request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		if len(_request.ASIN) == 0 {
			return kit.HTTPErrInvalidRequest
		}

		settings = AmazonCollectorSettings{
			ASIN: _request.ASIN,
		}
		jobdata = AmazonCollectorJobdata{
			LastDispatchedAt:    nil,
			LastDispatchedTasks: []string{},
			Cost:                0,
		}

		for _, _collector := range collectors {
			if _collector.Settings.(AmazonCollectorSettings).ASIN == _request.ASIN {
				collector = util.Pointer(_collector)
				break
			}
		}

	case CollectorTypeIAgora:
		var _request CollectorEndpointsPostIAgoraCollectorRequest
		err := json.Unmarshal(requestRaw, &_request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		if len(_request.Institution) == 0 {
			return kit.HTTPErrInvalidRequest
		}

		settings = IAgoraCollectorSettings{
			Institution: _request.Institution,
		}
		jobdata = IAgoraCollectorJobdata{
			LastCollectedAt: nil,
		}

		for _, _collector := range collectors {
			if _collector.Settings.(IAgoraCollectorSettings).Institution == _request.Institution {
				collector = util.Pointer(_collector)
				break
			}
		}

	case CollectorTypeWebhook:
		var _request CollectorEndpointsPostWebhookCollectorRequest
		err := json.Unmarshal(requestRaw, &_request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		settings = WebhookCollectorSettings{
			APIKey: util.RandomString(WEBHOOK_COLLECTOR_CALLBACK_SECRET_LENGTH),
		}
		jobdata = WebhookCollectorJobdata{}

		collector = nil

	case CollectorTypeWidget:
		var _request CollectorEndpointsPostWidgetCollectorRequest
		err := json.Unmarshal(requestRaw, &_request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		settings = WidgetCollectorSettings{
			ClientKey: util.RandomString(WIDGET_COLLECTOR_CALLBACK_SECRET_LENGTH),
		}
		jobdata = WidgetCollectorJobdata{}

		collector = nil

	default:
		return kit.HTTPErrServerGeneric
	}

	if collector != nil {
		if collector.DeletedAt == nil {
			return kit.HTTPErrInvalidRequest
		}

		collector.Settings = settings
		collector.DeletedAt = nil

		err := self.collectorRepository.UpdateSettingsAndDeleted(requestCtx, *collector)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}
	} else {
		collector = NewCollector()
		collector.ID = xid.New().String()
		collector.ProductID = requestProduct.ID
		collector.Type = request.Type
		collector.Settings = settings
		collector.Jobdata = jobdata
		collector.CreatedAt = time.Now()
		collector.DeletedAt = nil

		collector, err = self.collectorRepository.Create(requestCtx, *collector)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}
	}

	switch collector.Type {
	case CollectorTypeTrustpilot:
		err = self.enqueuer.Enqueue(requestCtx, TrustpilotCollectorDispatch, TrustpilotCollectorDispatchParams{
			CollectorID: collector.ID,
		}, asynq.MaxRetry(2), asynq.Unique(24*time.Hour))
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

	case CollectorTypePlayStore:
		err = self.enqueuer.Enqueue(requestCtx, PlayStoreCollectorDispatch, PlayStoreCollectorDispatchParams{
			CollectorID: collector.ID,
		}, asynq.MaxRetry(2), asynq.Unique(24*time.Hour))
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

	case CollectorTypeAppStore:
		err = self.enqueuer.Enqueue(requestCtx, AppStoreCollectorDispatch, AppStoreCollectorDispatchParams{
			CollectorID: collector.ID,
		}, asynq.MaxRetry(2), asynq.Unique(24*time.Hour))
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

	case CollectorTypeAmazon:
		err = self.enqueuer.Enqueue(requestCtx, AmazonCollectorDispatch, AmazonCollectorDispatchParams{
			CollectorID: collector.ID,
		}, asynq.MaxRetry(2), asynq.Unique(24*time.Hour))
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

	case CollectorTypeIAgora:
		err = self.enqueuer.Enqueue(requestCtx, IAgoraCollectorCollect, IAgoraCollectorCollectParams{
			CollectorID: collector.ID,
		}, asynq.MaxRetry(2), asynq.Unique(24*time.Hour))
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

	case CollectorTypeWebhook:

	case CollectorTypeWidget:

	default:
		return kit.HTTPErrServerGeneric
	}

	response := CollectorEndpointsPostCollectorResponse{}
	response.CollectorPayload = *NewCollectorPayload(*collector)

	return ctx.JSON(http.StatusOK, &response)
}

type CollectorEndpointsGetCollectorResponse struct {
	CollectorPayload
}

func (self *CollectorEndpoints) GetCollector(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestCollector := RequestCollector(requestCtx)

	response := CollectorEndpointsGetCollectorResponse{}
	response.CollectorPayload = *NewCollectorPayload(*requestCollector)

	return ctx.JSON(http.StatusOK, &response)
}

type CollectorEndpointsPutTrustpilotCollectorRequest struct {
	CollectorEndpointsPutCollectorRequest
}

type CollectorEndpointsPutPlayStoreCollectorRequest struct {
	CollectorEndpointsPutCollectorRequest
}

type CollectorEndpointsPutAppStoreCollectorRequest struct {
	CollectorEndpointsPutCollectorRequest
}

type CollectorEndpointsPutAmazonCollectorRequest struct {
	CollectorEndpointsPutCollectorRequest
}

type CollectorEndpointsPutIAgoraCollectorRequest struct {
	CollectorEndpointsPutCollectorRequest
}

type CollectorEndpointsPutWebhookCollectorRequest struct {
	CollectorEndpointsPutCollectorRequest
}

type CollectorEndpointsPutWidgetCollectorRequest struct {
	CollectorEndpointsPutCollectorRequest
}

type CollectorEndpointsPutCollectorRequest struct {
}

type CollectorEndpointsPutCollectorResponse struct {
	CollectorPayload
}

func (self *CollectorEndpoints) PutCollector(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestCollector := RequestCollector(requestCtx)

	switch requestCollector.Type {
	case CollectorTypeTrustpilot:
		request := CollectorEndpointsPutTrustpilotCollectorRequest{}

		err := ctx.Bind(&request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		settings := requestCollector.Settings.(TrustpilotCollectorSettings)

		requestCollector.Settings = settings

	case CollectorTypePlayStore:
		request := CollectorEndpointsPutPlayStoreCollectorRequest{}

		err := ctx.Bind(&request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		settings := requestCollector.Settings.(PlayStoreCollectorSettings)

		requestCollector.Settings = settings

	case CollectorTypeAppStore:
		request := CollectorEndpointsPutAppStoreCollectorRequest{}

		err := ctx.Bind(&request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		settings := requestCollector.Settings.(AppStoreCollectorSettings)

		requestCollector.Settings = settings

	case CollectorTypeAmazon:
		request := CollectorEndpointsPutAmazonCollectorRequest{}

		err := ctx.Bind(&request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		settings := requestCollector.Settings.(AmazonCollectorSettings)

		requestCollector.Settings = settings

	case CollectorTypeIAgora:
		request := CollectorEndpointsPutIAgoraCollectorRequest{}

		err := ctx.Bind(&request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		settings := requestCollector.Settings.(IAgoraCollectorSettings)

		requestCollector.Settings = settings

	case CollectorTypeWebhook:
		request := CollectorEndpointsPutWebhookCollectorRequest{}

		err := ctx.Bind(&request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		settings := requestCollector.Settings.(WebhookCollectorSettings)

		requestCollector.Settings = settings

	case CollectorTypeWidget:
		request := CollectorEndpointsPutWidgetCollectorRequest{}

		err := ctx.Bind(&request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		settings := requestCollector.Settings.(WidgetCollectorSettings)

		requestCollector.Settings = settings

	default:
		return kit.HTTPErrServerGeneric
	}

	err := self.collectorRepository.UpdateSettings(requestCtx, *requestCollector)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := CollectorEndpointsPutCollectorResponse{}
	response.CollectorPayload = *NewCollectorPayload(*requestCollector)

	return ctx.JSON(http.StatusOK, &response)
}

func (self *CollectorEndpoints) DeleteCollector(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestCollector := RequestCollector(requestCtx)

	err := self.collectorRepository.UpdateDeletedAt(requestCtx, requestCollector.ID, time.Now())
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, struct{}{})
}
