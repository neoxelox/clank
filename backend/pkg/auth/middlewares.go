package auth

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
	"backend/pkg/organization"
	"backend/pkg/user"
)

var (
	KeyAuthSignIn     kit.Key = kit.KeyBase + "auth:signin:"
	KeyRequestSession kit.Key = kit.KeyBase + "request:session"
)

func RequestSession(ctx context.Context) *Session {
	return ctx.Value(KeyRequestSession).(*Session) // nolint:forcetypeassert,errcheck
}

type AuthMiddlewares struct {
	config            config.Config
	observer          *kit.Observer
	authVerifier      *AuthVerifier
	sessionRepository SessionRepository
}

func NewAuthMiddlewares(observer *kit.Observer, authVerifier *AuthVerifier,
	sessionRepository SessionRepository, config config.Config) *AuthMiddlewares {
	return &AuthMiddlewares{
		config:            config,
		observer:          observer,
		authVerifier:      authVerifier,
		sessionRepository: sessionRepository,
	}
}

func (self *AuthMiddlewares) HandleToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestCtx := ctx.Request().Context()

		token, err := ctx.Cookie(TOKEN_COOKIE_NAME)
		if err != nil {
			return kit.HTTPErrUnauthorized.Cause(err)
		}

		result, err := self.authVerifier.CheckToken(requestCtx, AuthVerifierCheckTokenParams{
			Token: token.Value,
		})
		if err != nil {
			if ErrAuthVerifierGeneric.Is(err) {
				return kit.HTTPErrServerGeneric.Cause(err)
			}

			return kit.HTTPErrUnauthorized.Cause(err)
		}

		requestIP := ctx.RealIP()
		requestDevice := ctx.Request().UserAgent()
		newLocation := true
		for _, location := range result.Session.Metadata.Locations {
			if location.IP == requestIP && location.Device == requestDevice {
				newLocation = false
				break
			}
		}

		if newLocation {
			result.Session.Metadata.Locations = append(result.Session.Metadata.Locations, SessionMetadataLocation{
				IP:     requestIP,
				Device: requestDevice,
			})
		}
		result.Session.LastSeenAt = time.Now()

		err = self.sessionRepository.UpdateSeen(requestCtx, result.Session)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		ctx.SetRequest(ctx.Request().WithContext(context.WithValue(context.WithValue(context.WithValue(requestCtx,
			organization.KeyRequestOrganization, &result.Organization),
			KeyRequestSession, &result.Session),
			user.KeyRequestMe, &result.User)))

		return next(ctx)
	}
}

func (self *AuthMiddlewares) HandleRights(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestCtx := ctx.Request().Context()
		requestUser := user.RequestMe(requestCtx)

		if requestUser.Role != user.UserRoleAdmin {
			return kit.HTTPErrUnauthorized
		}

		return next(ctx)
	}
}
