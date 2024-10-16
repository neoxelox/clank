package organization

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
	"backend/pkg/util"
)

type OrganizationEndpoints struct {
	config                 config.Config
	observer               *kit.Observer
	organizationRepository OrganizationRepository
}

func NewOrganizationEndpoints(observer *kit.Observer, organizationRepository OrganizationRepository,
	config config.Config) *OrganizationEndpoints {
	return &OrganizationEndpoints{
		config:                 config,
		observer:               observer,
		organizationRepository: organizationRepository,
	}
}

type OrganizationEndpointsGetOrganizationResponse struct {
	OrganizationPayload
}

func (self *OrganizationEndpoints) GetOrganization(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestOrganization := RequestOrganization(requestCtx)

	response := OrganizationEndpointsGetOrganizationResponse{}
	response.OrganizationPayload = *NewOrganizationPayload(*requestOrganization)

	return ctx.JSON(http.StatusOK, &response)
}

type OrganizationEndpointsPutOrganizationRequest struct {
	Name    *string `json:"name"`
	Picture *string `json:"picture"`
}

type OrganizationEndpointsPutOrganizationResponse struct {
	OrganizationPayload
}

func (self *OrganizationEndpoints) PutOrganization(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestOrganization := RequestOrganization(requestCtx)

	request := OrganizationEndpointsPutOrganizationRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.Name != nil {
		if len(*request.Name) == 0 {
			return kit.HTTPErrInvalidRequest
		}

		requestOrganization.Name = *request.Name
	}

	if request.Picture != nil {
		if !util.SameOrigin(*request.Picture, self.config.CDN.BaseURL) {
			return kit.HTTPErrInvalidRequest
		}

		requestOrganization.Picture = *request.Picture
	}

	err = self.organizationRepository.UpdateProfile(requestCtx, *requestOrganization)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := OrganizationEndpointsPutOrganizationResponse{}
	response.OrganizationPayload = *NewOrganizationPayload(*requestOrganization)

	return ctx.JSON(http.StatusOK, &response)
}

func (self *OrganizationEndpoints) DeleteOrganization(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestOrganization := RequestOrganization(requestCtx)

	err := self.organizationRepository.Delete(requestCtx, requestOrganization.ID)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, struct{}{})
}

type OrganizationEndpointsGetOrganizationSettingsResponse struct {
	OrganizationPayloadSettings
}

func (self *OrganizationEndpoints) GetOrganizationSettings(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestOrganization := RequestOrganization(requestCtx)

	response := OrganizationEndpointsGetOrganizationSettingsResponse{}
	response.OrganizationPayloadSettings = NewOrganizationPayload(*requestOrganization).Settings

	return ctx.JSON(http.StatusOK, &response)
}

type OrganizationEndpointsPutOrganizationSettingsRequest struct {
	DomainSignIn *bool `json:"domain_sign_in"`
}

type OrganizationEndpointsPutOrganizationSettingsResponse struct {
	OrganizationPayloadSettings
}

func (self *OrganizationEndpoints) PutOrganizationSettings(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestOrganization := RequestOrganization(requestCtx)

	request := OrganizationEndpointsPutOrganizationSettingsRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.DomainSignIn != nil {
		if *request.DomainSignIn && !IsDomainSignInSupported(requestOrganization.Domain) {
			return kit.HTTPErrInvalidRequest
		}

		requestOrganization.Settings.DomainSignIn = *request.DomainSignIn
	}

	err = self.organizationRepository.UpdateSettings(requestCtx, *requestOrganization)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := OrganizationEndpointsPutOrganizationSettingsResponse{}
	response.OrganizationPayloadSettings = NewOrganizationPayload(*requestOrganization).Settings

	return ctx.JSON(http.StatusOK, &response)
}

type OrganizationEndpointsGetOrganizationUsageResponse struct {
	Capacity OrganizationPayloadCapacity `json:"capacity"`
	Usage    int                         `json:"usage"`
}

func (self *OrganizationEndpoints) GetOrganizationUsage(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestOrganization := RequestOrganization(requestCtx)

	organization := NewOrganizationPayload(*requestOrganization)

	response := OrganizationEndpointsGetOrganizationUsageResponse{}
	response.Capacity = organization.Capacity
	response.Usage = organization.Usage

	return ctx.JSON(http.StatusOK, &response)
}
