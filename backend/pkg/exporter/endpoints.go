package exporter

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"
	"github.com/neoxelox/kit/util"
	"github.com/rs/xid"

	"backend/pkg/config"
	"backend/pkg/organization"
	"backend/pkg/product"
)

type ExporterEndpoints struct {
	config             config.Config
	observer           *kit.Observer
	exporterRepository *ExporterRepository
}

func NewExporterEndpoints(observer *kit.Observer, exporterRepository *ExporterRepository,
	config config.Config) *ExporterEndpoints {
	return &ExporterEndpoints{
		config:             config,
		observer:           observer,
		exporterRepository: exporterRepository,
	}
}

type ExporterEndpointsListExportersResponse struct {
	Exporters []ExporterPayload `json:"exporters"`
}

func (self *ExporterEndpoints) ListExporters(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := product.RequestProduct(requestCtx)

	exporters, err := self.exporterRepository.ListByProductID(requestCtx, requestProduct.ID)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := ExporterEndpointsListExportersResponse{}
	response.Exporters = make([]ExporterPayload, 0, len(exporters))
	for _, exporter := range exporters {
		if exporter.DeletedAt != nil {
			continue
		}

		response.Exporters = append(response.Exporters, *NewExporterPayload(exporter))
	}

	return ctx.JSON(http.StatusOK, &response)
}

type ExporterEndpointsPostSlackExporterRequest struct {
	Channel string `json:"channel"`
}

type ExporterEndpointsPostJiraExporterRequest struct {
	Board string `json:"board"`
}

type ExporterEndpointsPostExporterRequest struct {
	Type string `json:"type"`
}

type ExporterEndpointsPostExporterResponse struct {
	ExporterPayload
}

func (self *ExporterEndpoints) PostExporter(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestOrganization := organization.RequestOrganization(requestCtx)
	requestProduct := product.RequestProduct(requestCtx)
	var requestRaw json.RawMessage
	request := ExporterEndpointsPostExporterRequest{}

	err := ctx.Bind(&requestRaw)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	err = json.Unmarshal(requestRaw, &request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if !IsExporterType(request.Type) {
		return kit.HTTPErrInvalidRequest
	}

	// TODO: Create an organization middleware that limits the usage of features by plan
	if requestOrganization.Plan != organization.OrganizationPlanEnterprise {
		return kit.HTTPErrInvalidRequest
	}

	exporters, err := self.exporterRepository.ListByProductIDAndType(requestCtx, requestProduct.ID, request.Type)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	var exporter *Exporter
	var settings any
	var jobdata any
	switch request.Type {
	case ExporterTypeSlack:
		var _request ExporterEndpointsPostSlackExporterRequest
		err := json.Unmarshal(requestRaw, &_request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		if len(_request.Channel) == 0 {
			return kit.HTTPErrInvalidRequest
		}

		settings = SlackExporterSettings{
			Channel: _request.Channel,
		}
		jobdata = SlackExporterJobdata{}

		for _, _exporter := range exporters {
			if _exporter.Settings.(SlackExporterSettings).Channel == _request.Channel {
				exporter = util.Pointer(_exporter)
				break
			}
		}

	case ExporterTypeJira:
		var _request ExporterEndpointsPostJiraExporterRequest
		err := json.Unmarshal(requestRaw, &_request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		if len(_request.Board) == 0 {
			return kit.HTTPErrInvalidRequest
		}

		settings = JiraExporterSettings{
			Board: _request.Board,
		}
		jobdata = JiraExporterJobdata{}

		for _, _exporter := range exporters {
			if _exporter.Settings.(JiraExporterSettings).Board == _request.Board {
				exporter = util.Pointer(_exporter)
				break
			}
		}

	default:
		return kit.HTTPErrServerGeneric
	}

	if exporter != nil {
		if exporter.DeletedAt == nil {
			return kit.HTTPErrInvalidRequest
		}

		exporter.Settings = settings
		exporter.DeletedAt = nil

		err := self.exporterRepository.UpdateSettingsAndDeleted(requestCtx, *exporter)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}
	} else {
		exporter = NewExporter()
		exporter.ID = xid.New().String()
		exporter.ProductID = requestProduct.ID
		exporter.Type = request.Type
		exporter.Settings = settings
		exporter.Jobdata = jobdata
		exporter.CreatedAt = time.Now()
		exporter.DeletedAt = nil

		exporter, err = self.exporterRepository.Create(requestCtx, *exporter)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}
	}

	switch exporter.Type {
	case ExporterTypeSlack:

	case ExporterTypeJira:

	default:
		return kit.HTTPErrServerGeneric
	}

	response := ExporterEndpointsPostExporterResponse{}
	response.ExporterPayload = *NewExporterPayload(*exporter)

	return ctx.JSON(http.StatusOK, &response)
}

type ExporterEndpointsGetExporterResponse struct {
	ExporterPayload
}

func (self *ExporterEndpoints) GetExporter(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestExporter := RequestExporter(requestCtx)

	response := ExporterEndpointsGetExporterResponse{}
	response.ExporterPayload = *NewExporterPayload(*requestExporter)

	return ctx.JSON(http.StatusOK, &response)
}

type ExporterEndpointsPutSlackExporterRequest struct {
	ExporterEndpointsPutExporterRequest
}

type ExporterEndpointsPutJiraExporterRequest struct {
	ExporterEndpointsPutExporterRequest
}

type ExporterEndpointsPutExporterRequest struct {
}

type ExporterEndpointsPutExporterResponse struct {
	ExporterPayload
}

func (self *ExporterEndpoints) PutExporter(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestExporter := RequestExporter(requestCtx)

	switch requestExporter.Type {
	case ExporterTypeSlack:
		request := ExporterEndpointsPutSlackExporterRequest{}

		err := ctx.Bind(&request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		settings := requestExporter.Settings.(SlackExporterSettings)

		requestExporter.Settings = settings

	case ExporterTypeJira:
		request := ExporterEndpointsPutJiraExporterRequest{}

		err := ctx.Bind(&request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		settings := requestExporter.Settings.(JiraExporterSettings)

		requestExporter.Settings = settings

	default:
		return kit.HTTPErrServerGeneric
	}

	err := self.exporterRepository.UpdateSettings(requestCtx, *requestExporter)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := ExporterEndpointsPutExporterResponse{}
	response.ExporterPayload = *NewExporterPayload(*requestExporter)

	return ctx.JSON(http.StatusOK, &response)
}

func (self *ExporterEndpoints) DeleteExporter(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestExporter := RequestExporter(requestCtx)

	err := self.exporterRepository.UpdateDeletedAt(requestCtx, requestExporter.ID, time.Now())
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, struct{}{})
}
