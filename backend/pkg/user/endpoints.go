package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"
	"github.com/rs/xid"

	"backend/pkg/brevo"
	"backend/pkg/config"
	"backend/pkg/organization"
	"backend/pkg/util"

	"github.com/badoux/checkmail"
	emailproviders "gomodules.xyz/email-providers"
)

const (
	USER_ENDPOINTS_INVITATION_EMAIL_SUBJECT  = "%s has invited you to join the %s organization"
	USER_ENDPOINTS_INVITATION_EMAIL_TEMPLATE = "emails/invitation.html"
)

type UserEndpoints struct {
	config                 config.Config
	observer               *kit.Observer
	database               *kit.Database
	renderer               *kit.Renderer
	brevoService           *brevo.BrevoService
	userRepository         UserRepository
	invitationRepository   InvitationRepository
	organizationRepository organization.OrganizationRepository
}

func NewUserEndpoints(observer *kit.Observer, database *kit.Database, renderer *kit.Renderer,
	brevoService *brevo.BrevoService, userRepository UserRepository, invitationRepository InvitationRepository,
	organizationRepository organization.OrganizationRepository, config config.Config) *UserEndpoints {
	return &UserEndpoints{
		config:                 config,
		observer:               observer,
		database:               database,
		renderer:               renderer,
		brevoService:           brevoService,
		userRepository:         userRepository,
		invitationRepository:   invitationRepository,
		organizationRepository: organizationRepository,
	}
}

type UserEndpointsGetMeResponse struct {
	UserPayload
}

func (self *UserEndpoints) GetMe(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestUser := RequestMe(requestCtx)

	response := UserEndpointsGetMeResponse{}
	response.UserPayload = *NewUserPayload(*requestUser)

	return ctx.JSON(http.StatusOK, &response)
}

type UserEndpointsPutMeRequest struct {
	Name    *string `json:"name"`
	Picture *string `json:"picture"`
}

type UserEndpointsPutMeResponse struct {
	UserPayload
}

func (self *UserEndpoints) PutMe(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestUser := RequestMe(requestCtx)
	request := UserEndpointsPutMeRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.Name != nil {
		if len(*request.Name) == 0 {
			return kit.HTTPErrInvalidRequest
		}

		requestUser.Name = *request.Name
	}

	if request.Picture != nil {
		if !util.SameOrigin(*request.Picture, self.config.CDN.BaseURL) {
			return kit.HTTPErrInvalidRequest
		}

		requestUser.Picture = *request.Picture
	}

	err = self.userRepository.UpdateProfile(requestCtx, *requestUser)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := UserEndpointsPutMeResponse{}
	response.UserPayload = *NewUserPayload(*requestUser)

	return ctx.JSON(http.StatusOK, &response)
}

func (self *UserEndpoints) DeleteMe(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestUser := RequestMe(requestCtx)
	requestOrganization := organization.RequestOrganization(requestCtx)

	count, err := self.userRepository.CountNotDeletedByOrganizationID(requestCtx, requestOrganization.ID)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	if count <= 1 {
		err := self.organizationRepository.Delete(requestCtx, requestOrganization.ID)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}
	} else {
		// TODO: If there are no admins left... make someone random an admin?
		err := self.userRepository.UpdateDeletedAt(requestCtx, requestUser.ID, time.Now())
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}
	}

	return ctx.JSON(http.StatusOK, struct{}{})
}

type UserEndpointsGetMySettingsResponse struct {
	UserPayloadSettings
}

func (self *UserEndpoints) GetMySettings(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestUser := RequestMe(requestCtx)

	response := UserEndpointsGetMySettingsResponse{}
	response.UserPayloadSettings = NewUserPayload(*requestUser).Settings

	return ctx.JSON(http.StatusOK, &response)
}

type UserEndpointsPutMySettingsRequest struct {
}

type UserEndpointsPutMySettingsResponse struct {
	UserPayloadSettings
}

func (self *UserEndpoints) PutMySettings(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestUser := RequestMe(requestCtx)
	request := UserEndpointsPutMySettingsRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	err = self.userRepository.UpdateSettings(requestCtx, *requestUser)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := UserEndpointsPutMySettingsResponse{}
	response.UserPayloadSettings = NewUserPayload(*requestUser).Settings

	return ctx.JSON(http.StatusOK, &response)
}

type UserEndpointsListUsersResponse struct {
	Users []UserPayload `json:"users"`
}

func (self *UserEndpoints) ListUsers(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestOrganization := organization.RequestOrganization(requestCtx)

	users, err := self.userRepository.ListByOrganizationID(requestCtx, requestOrganization.ID)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := UserEndpointsListUsersResponse{}
	response.Users = make([]UserPayload, 0, len(users))
	for _, user := range users {
		response.Users = append(response.Users, *NewUserPayload(user))
	}

	return ctx.JSON(http.StatusOK, &response)
}

type UserEndpointsGetUserResponse struct {
	UserPayload
}

func (self *UserEndpoints) GetUser(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestUser := RequestUser(requestCtx)

	response := UserEndpointsGetUserResponse{}
	response.UserPayload = *NewUserPayload(*requestUser)

	return ctx.JSON(http.StatusOK, &response)
}

type UserEndpointsPutUserRequest struct {
	Role *string `json:"role"`
}

type UserEndpointsPutUserResponse struct {
	UserPayload
}

func (self *UserEndpoints) PutUser(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestUser := RequestUser(requestCtx)
	request := UserEndpointsPutUserRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.Role != nil {
		if !IsUserRole(*request.Role) {
			return kit.HTTPErrInvalidRequest
		}

		requestUser.Role = *request.Role
	}

	err = self.userRepository.UpdateRole(requestCtx, requestUser.ID, requestUser.Role)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := UserEndpointsPutUserResponse{}
	response.UserPayload = *NewUserPayload(*requestUser)

	return ctx.JSON(http.StatusOK, &response)
}

func (self *UserEndpoints) DeleteUser(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestUser := RequestUser(requestCtx)

	err := self.userRepository.UpdateDeletedAt(requestCtx, requestUser.ID, time.Now())
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, struct{}{})
}

type UserEndpointsListInvitationsResponse struct {
	Invitations []InvitationPayload `json:"invitations"`
}

func (self *UserEndpoints) ListInvitations(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestOrganization := organization.RequestOrganization(requestCtx)

	invitations, err := self.invitationRepository.ListByOrganizationID(requestCtx, requestOrganization.ID)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := UserEndpointsListInvitationsResponse{}
	response.Invitations = make([]InvitationPayload, 0, len(invitations))
	for _, invitation := range invitations {
		response.Invitations = append(response.Invitations, *NewInvitationPayload(invitation))
	}

	return ctx.JSON(http.StatusOK, &response)
}

type UserEndpointsPostInvitationRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserEndpointsPostInvitationResponse struct {
	InvitationPayload
}

func (self *UserEndpoints) PostInvitation(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestUser := RequestMe(requestCtx)
	requestOrganization := organization.RequestOrganization(requestCtx)
	request := UserEndpointsPostInvitationRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	// TODO: Create an organization middleware that limits the usage of features by plan
	if requestOrganization.Plan == organization.OrganizationPlanDemo {
		return kit.HTTPErrInvalidRequest
	}

	// TODO: Create an organization middleware that limits the usage of features by plan
	if requestOrganization.Plan == organization.OrganizationPlanTrial {
		userCount, err := self.userRepository.CountNotDeletedByOrganizationID(requestCtx, requestOrganization.ID)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		if userCount >= organization.ORGANIZATION_PLAN_TRIAL_MAX_MEMBERS {
			return kit.HTTPErrInvalidRequest
		}

		invitationCount, err := self.invitationRepository.CountNotExpiredByOrganizationID(requestCtx, requestOrganization.ID)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		if invitationCount >= (organization.ORGANIZATION_PLAN_TRIAL_MAX_MEMBERS - userCount) {
			return kit.HTTPErrInvalidRequest
		}
	}

	err = checkmail.ValidateFormat(request.Email)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if emailproviders.IsDisposableEmail(request.Email) {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if !IsUserRole(request.Role) {
		return kit.HTTPErrInvalidRequest
	}

	user, err := self.userRepository.GetByEmail(requestCtx, request.Email)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	if user != nil && user.DeletedAt == nil {
		return kit.HTTPErrInvalidRequest
	}

	invitation, err := self.invitationRepository.GetByEmail(requestCtx, request.Email)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	if invitation != nil {
		if time.Now().Before(invitation.ExpiresAt) {
			return kit.HTTPErrInvalidRequest
		}

		err = self.invitationRepository.DeleteByID(requestCtx, invitation.ID)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}
	}

	invitation = NewInvitation()
	invitation.ID = xid.New().String()
	invitation.OrganizationID = requestOrganization.ID
	invitation.Email = request.Email
	invitation.Role = request.Role
	invitation.ExpiresAt = time.Now().Add(INVITATION_EXPIRATION)

	err = self.database.Transaction(requestCtx, nil, func(requestCtx context.Context) error {
		invitation, err = self.invitationRepository.Create(requestCtx, *invitation)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		body, err := self.renderer.RenderString(USER_ENDPOINTS_INVITATION_EMAIL_TEMPLATE,
			map[string]any{"Inviter": requestUser.Name, "Organization": requestOrganization.Name})
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		err = self.brevoService.SendEmail(requestCtx, brevo.BrevoServiceSendEmailParams{
			Receivers: []string{invitation.Email},
			Subject: fmt.Sprintf(USER_ENDPOINTS_INVITATION_EMAIL_SUBJECT,
				requestUser.Name, requestOrganization.Name),
			Body: body,
		})
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	response := UserEndpointsPostInvitationResponse{}
	response.InvitationPayload = *NewInvitationPayload(*invitation)

	return ctx.JSON(http.StatusOK, &response)
}

type UserEndpointsGetInvitationResponse struct {
	InvitationPayload
}

func (self *UserEndpoints) GetInvitation(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestInvitation := RequestInvitation(requestCtx)

	response := UserEndpointsGetInvitationResponse{}
	response.InvitationPayload = *NewInvitationPayload(*requestInvitation)

	return ctx.JSON(http.StatusOK, &response)
}

func (self *UserEndpoints) DeleteInvitation(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestInvitation := RequestInvitation(requestCtx)

	err := self.invitationRepository.DeleteByID(requestCtx, requestInvitation.ID)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, struct{}{})
}
