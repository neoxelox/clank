package user

import (
	"context"
	"net/http"

	"backend/pkg/config"
	"backend/pkg/organization"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"
)

var (
	KeyRequestMe         kit.Key = kit.KeyBase + "request:me"
	KeyRequestUser       kit.Key = kit.KeyBase + "request:user"
	KeyRequestInvitation kit.Key = kit.KeyBase + "request:invitation"
)

func RequestMe(ctx context.Context) *User {
	return ctx.Value(KeyRequestMe).(*User) // nolint:forcetypeassert,errcheck
}

func RequestUser(ctx context.Context) *User {
	return ctx.Value(KeyRequestUser).(*User) // nolint:forcetypeassert,errcheck
}

func RequestInvitation(ctx context.Context) *Invitation {
	return ctx.Value(KeyRequestInvitation).(*Invitation) // nolint:forcetypeassert,errcheck
}

type UserMiddlewares struct {
	config               config.Config
	observer             *kit.Observer
	userRepository       UserRepository
	invitationRepository InvitationRepository
}

func NewUserMiddlewares(observer *kit.Observer, userRepository UserRepository,
	invitationRepository InvitationRepository, config config.Config) *UserMiddlewares {
	return &UserMiddlewares{
		config:               config,
		observer:             observer,
		userRepository:       userRepository,
		invitationRepository: invitationRepository,
	}
}

func (self *UserMiddlewares) HandleUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestCtx := ctx.Request().Context()
		requestMe := RequestMe(requestCtx)
		requestOrganization := organization.RequestOrganization(requestCtx)

		user, err := self.userRepository.GetByID(requestCtx, ctx.Param("user_id"))
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		if user == nil {
			return kit.HTTPErrInvalidRequest
		}

		if user.ID == requestMe.ID {
			return kit.HTTPErrInvalidRequest
		}

		if user.OrganizationID != requestOrganization.ID {
			return kit.HTTPErrUnauthorized
		}

		if ctx.Request().Method != http.MethodGet {
			if user.DeletedAt != nil {
				return kit.HTTPErrUnauthorized
			}
		}

		ctx.SetRequest(ctx.Request().WithContext(context.WithValue(requestCtx, KeyRequestUser, user)))

		return next(ctx)
	}
}

func (self *UserMiddlewares) HandleInvitation(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestCtx := ctx.Request().Context()
		requestOrganization := organization.RequestOrganization(requestCtx)

		invitation, err := self.invitationRepository.GetByID(requestCtx, ctx.Param("invitation_id"))
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		if invitation == nil {
			return kit.HTTPErrInvalidRequest
		}

		if invitation.OrganizationID != requestOrganization.ID {
			return kit.HTTPErrUnauthorized
		}

		ctx.SetRequest(ctx.Request().WithContext(context.WithValue(requestCtx, KeyRequestInvitation, invitation)))

		return next(ctx)
	}
}
