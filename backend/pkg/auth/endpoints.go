package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
)

type AuthEndpoints struct {
	config        config.Config
	observer      *kit.Observer
	authProcessor *AuthProcessor
}

func NewAuthEndpoints(observer *kit.Observer, authProcessor *AuthProcessor,
	config config.Config) *AuthEndpoints {
	return &AuthEndpoints{
		config:        config,
		observer:      observer,
		authProcessor: authProcessor,
	}
}

func (self *AuthEndpoints) newStateCookie(value string) *http.Cookie {
	maxAge := int(STATE_EXPIRATION.Seconds())
	if len(value) == 0 {
		maxAge = -1
	}

	return &http.Cookie{
		Name:     STATE_COOKIE_NAME,
		Value:    value,
		Path:     "/",
		Domain:   self.config.Service.Domain,
		MaxAge:   maxAge,
		Secure:   self.config.Service.Environment == kit.EnvProduction,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func (self *AuthEndpoints) newTokenCookie(value string) *http.Cookie {
	maxAge := int(TOKEN_EXPIRATION.Seconds())
	if len(value) == 0 {
		maxAge = -1
	}

	return &http.Cookie{
		Name:     TOKEN_COOKIE_NAME,
		Value:    value,
		Path:     "/",
		Domain:   self.config.Service.Domain,
		MaxAge:   maxAge,
		Secure:   self.config.Service.Environment == kit.EnvProduction,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

type AuthEndpointsPostSignInStartRequest struct {
	RedirectTo *string `json:"redirect_to"`
}

type AuthEndpointsPostSignInEmailStartRequest struct {
	AuthEndpointsPostSignInStartRequest
	Email string `json:"email"`
}

type AuthEndpointsPostSignInEmailStartResponse struct {
	SignInCodeState string `json:"sign_in_code_state"`
	SignInCodeID    string `json:"sign_in_code_id"`
}

type AuthEndpointsPostSignInOAuthStartRequest struct {
	AuthEndpointsPostSignInStartRequest
}

type AuthEndpointsPostSignInOAuthStartResponse struct {
	AuthURL string `json:"auth_url"`
}

type AuthEndpointsPostSignInSAMLStartRequest struct {
	AuthEndpointsPostSignInStartRequest
}

type AuthEndpointsPostSignInSAMLStartResponse struct {
}

func (self *AuthEndpoints) PostSignInStart(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()

	provider := strings.ToUpper(ctx.Param("provider"))
	if !IsSessionProvider(provider) {
		return kit.HTTPErrNotFound
	}

	var token *string
	tokenCookie, err := ctx.Cookie(TOKEN_COOKIE_NAME)
	if err == nil {
		token = &tokenCookie.Value
	}

	switch provider {
	case SessionProviderEmail:
		request := AuthEndpointsPostSignInEmailStartRequest{}

		err := ctx.Bind(&request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		result, err := self.authProcessor.StartEmailSignIn(requestCtx, AuthProcessorStartEmailSignInParams{
			AuthProcessorStartSignInParams: AuthProcessorStartSignInParams{
				RedirectTo:   request.RedirectTo,
				CurrentToken: token,
			},
			Email: request.Email,
		})
		switch {
		case err == nil:
			response := AuthEndpointsPostSignInEmailStartResponse{}
			response.SignInCodeState = result.State
			response.SignInCodeID = result.SignInCodeID

			ctx.SetCookie(self.newStateCookie(result.State))

			return ctx.JSON(http.StatusOK, &response)

		case ErrAuthProcessorInvalidEmail.Is(err), ErrAuthProcessorSignInCodeAlreadyRequested.Is(err):
			return kit.HTTPErrInvalidRequest.Cause(err)

		default:
			return kit.HTTPErrServerGeneric.Cause(err)
		}

	case SessionProviderSAML:
		request := AuthEndpointsPostSignInSAMLStartRequest{}

		err := ctx.Bind(&request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		result, err := self.authProcessor.StartSAMLSignIn(requestCtx, AuthProcessorStartSAMLSignInParams{
			AuthProcessorStartSignInParams: AuthProcessorStartSignInParams{
				RedirectTo:   request.RedirectTo,
				CurrentToken: token,
			},
		})
		switch {
		case err == nil:
			response := AuthEndpointsPostSignInSAMLStartResponse{}

			ctx.SetCookie(self.newStateCookie(result.State))

			return ctx.JSON(http.StatusOK, &response)

		default:
			return kit.HTTPErrServerGeneric.Cause(err)
		}

	default:
		request := AuthEndpointsPostSignInOAuthStartRequest{}

		err := ctx.Bind(&request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		result, err := self.authProcessor.StartOAuthSignIn(requestCtx, AuthProcessorStartOAuthSignInParams{
			AuthProcessorStartSignInParams: AuthProcessorStartSignInParams{
				RedirectTo:   request.RedirectTo,
				CurrentToken: token,
			},
			Provider: provider,
		})
		switch {
		case err == nil:
			response := AuthEndpointsPostSignInOAuthStartResponse{}
			response.AuthURL = result.AuthURL

			ctx.SetCookie(self.newStateCookie(result.State))

			return ctx.JSON(http.StatusOK, &response)

		case ErrAuthProcessorInvalidOAuthProvider.Is(err):
			return kit.HTTPErrInvalidRequest.Cause(err)

		default:
			return kit.HTTPErrServerGeneric.Cause(err)
		}
	}
}

type AuthEndpointsPostSignInEndRequest struct {
	State string `json:"state"`
}

type AuthEndpointsPostSignInEndResponse struct {
	RedirectTo *string `json:"redirect_to"`
}

type AuthEndpointsPostSignInEmailEndRequest struct {
	AuthEndpointsPostSignInEndRequest
	SignInCodeID   string `json:"sign_in_code_id"`
	SignInCodeCode string `json:"sign_in_code_code"`
}

type AuthEndpointsPostSignInEmailEndResponse struct {
	AuthEndpointsPostSignInEndResponse
}

type AuthEndpointsPostSignInOAuthEndRequest struct {
	AuthEndpointsPostSignInEndRequest
	AuthResult string `json:"auth_result"`
}

type AuthEndpointsPostSignInOAuthEndResponse struct {
	AuthEndpointsPostSignInEndResponse
}

type AuthEndpointsPostSignInSAMLEndRequest struct {
	AuthEndpointsPostSignInEndRequest
}

type AuthEndpointsPostSignInSAMLEndResponse struct {
	AuthEndpointsPostSignInEndResponse
}

func (self *AuthEndpoints) PostSignInEnd(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()

	provider := strings.ToUpper(ctx.Param("provider"))
	if !IsSessionProvider(provider) {
		return kit.HTTPErrNotFound
	}

	state, err := ctx.Cookie(STATE_COOKIE_NAME)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}
	ip := ctx.RealIP()
	device := ctx.Request().UserAgent()

	switch provider {
	case SessionProviderEmail:
		request := AuthEndpointsPostSignInEmailEndRequest{}

		err := ctx.Bind(&request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		result, err := self.authProcessor.EndEmailSignIn(requestCtx, AuthProcessorEndEmailSignInParams{
			AuthProcessorEndSignInParams: AuthProcessorEndSignInParams{
				StartState: state.Value,
				EndState:   request.State,
				IP:         ip,
				Device:     device,
			},
			SignInCodeID:   request.SignInCodeID,
			SignInCodeCode: request.SignInCodeCode,
		})
		switch {
		case err == nil:
			response := AuthEndpointsPostSignInEmailEndResponse{}
			response.RedirectTo = result.RedirectTo

			ctx.SetCookie(self.newStateCookie(""))
			ctx.SetCookie(self.newTokenCookie(result.Token))

			return ctx.JSON(http.StatusOK, &response)

		case ErrAuthProcessorInvalidSignInCode.Is(err), ErrAuthProcessorInvalidState.Is(err),
			ErrAuthProcessorUnauthorizedUser.Is(err):
			return kit.HTTPErrInvalidRequest.Cause(err)

		default:
			return kit.HTTPErrServerGeneric.Cause(err)
		}

	case SessionProviderSAML:
		request := AuthEndpointsPostSignInSAMLEndRequest{}

		err := ctx.Bind(&request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		result, err := self.authProcessor.EndSAMLSignIn(requestCtx, AuthProcessorEndSAMLSignInParams{
			AuthProcessorEndSignInParams: AuthProcessorEndSignInParams{
				StartState: state.Value,
				EndState:   request.State,
				IP:         ip,
				Device:     device,
			},
		})
		switch {
		case err == nil:
			response := AuthEndpointsPostSignInSAMLEndResponse{}
			response.RedirectTo = result.RedirectTo

			ctx.SetCookie(self.newStateCookie(""))
			ctx.SetCookie(self.newTokenCookie(result.Token))

			return ctx.JSON(http.StatusOK, &response)

		case ErrAuthProcessorInvalidState.Is(err), ErrAuthProcessorUnauthorizedUser.Is(err):
			return kit.HTTPErrInvalidRequest.Cause(err)

		default:
			return kit.HTTPErrServerGeneric.Cause(err)
		}

	default:
		request := AuthEndpointsPostSignInOAuthEndRequest{}

		err := ctx.Bind(&request)
		if err != nil {
			return kit.HTTPErrInvalidRequest.Cause(err)
		}

		result, err := self.authProcessor.EndOAuthSignIn(requestCtx, AuthProcessorEndOAuthSignInParams{
			AuthProcessorEndSignInParams: AuthProcessorEndSignInParams{
				StartState: state.Value,
				EndState:   request.State,
				IP:         ip,
				Device:     device,
			},
			Provider:   provider,
			AuthResult: request.AuthResult,
		})
		switch {
		case err == nil:
			response := AuthEndpointsPostSignInOAuthEndResponse{}
			response.RedirectTo = result.RedirectTo

			ctx.SetCookie(self.newStateCookie(""))
			ctx.SetCookie(self.newTokenCookie(result.Token))

			return ctx.JSON(http.StatusOK, &response)

		case ErrAuthProcessorInvalidOAuthProvider.Is(err), ErrAuthProcessorInvalidOAuthResult.Is(err),
			ErrAuthProcessorInvalidState.Is(err), ErrAuthProcessorUnauthorizedUser.Is(err):
			return kit.HTTPErrInvalidRequest.Cause(err)

		default:
			return kit.HTTPErrServerGeneric.Cause(err)
		}
	}
}

func (self *AuthEndpoints) PostSignOut(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestSession := RequestSession(requestCtx)

	err := self.authProcessor.SignOut(requestCtx, AuthProcessorSignOutParams{
		Session: *requestSession,
	})
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	ctx.SetCookie(self.newStateCookie(""))
	ctx.SetCookie(self.newTokenCookie(""))

	return ctx.JSON(http.StatusOK, struct{}{})
}
