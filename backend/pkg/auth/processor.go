package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/url"
	"strings"
	"time"

	"backend/pkg/brevo"
	"backend/pkg/config"
	"backend/pkg/organization"
	"backend/pkg/user"
	"backend/pkg/util"

	"github.com/badoux/checkmail"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/google"
	"github.com/neoxelox/errors"
	"github.com/neoxelox/kit"
	kitUtil "github.com/neoxelox/kit/util"
	"github.com/rs/xid"
	"github.com/vk-rv/pvx"
	emailproviders "gomodules.xyz/email-providers"
)

const (
	AUTH_PROCESSOR_DEVELOPMENT_SIGN_IN_CODE    = "123456"
	AUTH_PROCESSOR_SIGN_IN_CODE_EMAIL_SUBJECT  = "Sign in Code"
	AUTH_PROCESSOR_SIGN_IN_CODE_EMAIL_TEMPLATE = "emails/sign_in_code.html"
)

var (
	ErrAuthProcessorGeneric                    = errors.New("auth processor failed")
	ErrAuthProcessorInvalidEmail               = errors.New("email is invalid")
	ErrAuthProcessorSignInCodeAlreadyRequested = errors.New("sign in code already requested")
	ErrAuthProcessorInvalidOAuthProvider       = errors.New("oauth provider is not registered")
	ErrAuthProcessorInvalidState               = errors.New("state is invalid")
	ErrAuthProcessorInvalidSignInCode          = errors.New("sign in code is invalid")
	ErrAuthProcessorInvalidOAuthResult         = errors.New("oauth query/form result values are invalid")
	ErrAuthProcessorUnauthorizedUser           = errors.New("user is unauthorized")
)

type AuthProcessor struct {
	config                 config.Config
	observer               *kit.Observer
	cryptKey               []byte
	tokenizer              *pvx.ProtoV4Local
	database               *kit.Database
	signInCodeRepository   SignInCodeRepository
	renderer               *kit.Renderer
	brevoService           *brevo.BrevoService
	authVerifier           *AuthVerifier
	userRepository         user.UserRepository
	invitationRepository   user.InvitationRepository
	organizationRepository organization.OrganizationRepository
	sessionRepository      SessionRepository
}

func NewAuthProcessor(observer *kit.Observer, database *kit.Database, signInCodeRepository SignInCodeRepository,
	renderer *kit.Renderer, brevoService *brevo.BrevoService, authVerifier *AuthVerifier,
	userRepository user.UserRepository, invitationRepository user.InvitationRepository,
	organizationRepository organization.OrganizationRepository, sessionRepository SessionRepository,
	config config.Config) *AuthProcessor {
	cryptKey := sha256.Sum256([]byte(config.Auth.CryptKey))

	tokenizer := pvx.NewPV4Local()

	goth.UseProviders(
		google.New(config.Auth.GoogleID, config.Auth.GoogleSecret,
			config.Frontend.BaseURL+"/dash/signin/google/end", "email", "profile"),
		apple.New(config.Auth.AppleID, config.Auth.AppleSecret,
			config.Frontend.BaseURL+"/dash/signin/apple/end", nil, "email", "name"),
		amazon.New(config.Auth.AmazonID, config.Auth.AmazonSecret,
			config.Frontend.BaseURL+"/dash/signin/amazon/end", "profile"),
	)

	return &AuthProcessor{
		config:                 config,
		observer:               observer,
		cryptKey:               cryptKey[:],
		tokenizer:              tokenizer,
		database:               database,
		signInCodeRepository:   signInCodeRepository,
		renderer:               renderer,
		brevoService:           brevoService,
		authVerifier:           authVerifier,
		userRepository:         userRepository,
		invitationRepository:   invitationRepository,
		organizationRepository: organizationRepository,
		sessionRepository:      sessionRepository,
	}
}

type AuthProcessorStartSignInParams struct {
	RedirectTo   *string
	CurrentToken *string
}

type AuthProcessorStartSignInResult struct {
	State string
}

func (self *AuthProcessor) startSignIn(ctx context.Context,
	params AuthProcessorStartSignInParams) (*AuthProcessorStartSignInResult, error) {
	statePayload := NewState()
	statePayload.RedirectTo = params.RedirectTo
	statePayload.PreviousToken = params.CurrentToken
	statePayload.ExpiresAt = time.Now().Add(STATE_EXPIRATION)
	statePayload.Nonce = kitUtil.RandomString(8)

	state, err := self.tokenizer.Encrypt(pvx.NewSymmetricKey(self.cryptKey, pvx.Version4), statePayload)
	if err != nil {
		return nil, ErrAuthProcessorGeneric.Raise().Cause(err)
	}

	result := AuthProcessorStartSignInResult{}
	result.State = state

	return &result, nil
}

type AuthProcessorStartEmailSignInParams struct {
	AuthProcessorStartSignInParams
	Email string
}

type AuthProcessorStartEmailSignInResult struct {
	AuthProcessorStartSignInResult
	SignInCodeID string
}

func (self *AuthProcessor) StartEmailSignIn(ctx context.Context,
	params AuthProcessorStartEmailSignInParams) (*AuthProcessorStartEmailSignInResult, error) {
	startResult, err := self.startSignIn(ctx, params.AuthProcessorStartSignInParams)
	if err != nil {
		return nil, err
	}

	err = checkmail.ValidateFormat(params.Email)
	if err != nil {
		return nil, ErrAuthProcessorInvalidEmail.Raise().Cause(err)
	}

	if emailproviders.IsDisposableEmail(params.Email) {
		return nil, ErrAuthProcessorInvalidEmail.Raise().With("disposable email")
	}

	signInCode, err := self.signInCodeRepository.GetByEmail(ctx, params.Email)
	if err != nil {
		return nil, ErrAuthProcessorGeneric.Raise().Cause(err)
	}

	if signInCode != nil {
		if time.Now().Before(signInCode.ExpiresAt) {
			return nil, ErrAuthProcessorSignInCodeAlreadyRequested.Raise()
		}

		err = self.signInCodeRepository.DeleteByID(ctx, signInCode.ID)
		if err != nil {
			return nil, ErrAuthProcessorGeneric.Raise().Cause(err)
		}
	}

	signInCode = NewSignInCode()
	signInCode.ID = xid.New().String()
	signInCode.Email = params.Email
	signInCode.Code = strings.ToUpper(kitUtil.RandomString(SIGN_IN_CODE_LENGTH))
	if self.config.Service.Environment != kit.EnvProduction {
		signInCode.Code = AUTH_PROCESSOR_DEVELOPMENT_SIGN_IN_CODE
	}
	signInCode.Attempts = 0
	signInCode.ExpiresAt = time.Now().Add(SIGN_IN_CODE_EXPIRATION)

	signInCodeURL := fmt.Sprintf("%s/dash/signin/email/end?id=%s&code=%s&state=%s",
		self.config.Frontend.BaseURL, url.QueryEscape(signInCode.ID),
		url.QueryEscape(signInCode.Code), url.QueryEscape(startResult.State))

	err = self.database.Transaction(ctx, nil, func(ctx context.Context) error {
		signInCode, err = self.signInCodeRepository.Create(ctx, *signInCode)
		if err != nil {
			return ErrAuthProcessorGeneric.Raise().Cause(err)
		}

		body, err := self.renderer.RenderString(AUTH_PROCESSOR_SIGN_IN_CODE_EMAIL_TEMPLATE,
			map[string]any{"Code": signInCode.Code, "URL": signInCodeURL})
		if err != nil {
			return ErrAuthProcessorGeneric.Raise().Cause(err)
		}

		err = self.brevoService.SendEmail(ctx, brevo.BrevoServiceSendEmailParams{
			Receivers: []string{signInCode.Email},
			Subject:   AUTH_PROCESSOR_SIGN_IN_CODE_EMAIL_SUBJECT,
			Body:      body,
		})
		if err != nil {
			return ErrAuthProcessorGeneric.Raise().Cause(err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	result := AuthProcessorStartEmailSignInResult{}
	result.AuthProcessorStartSignInResult = *startResult
	result.SignInCodeID = signInCode.ID

	return &result, nil
}

type AuthProcessorStartOAuthSignInParams struct {
	AuthProcessorStartSignInParams
	Provider string
}

type AuthProcessorStartOAuthSignInResult struct {
	AuthProcessorStartSignInResult
	AuthURL string
}

func (self *AuthProcessor) StartOAuthSignIn(ctx context.Context,
	params AuthProcessorStartOAuthSignInParams) (*AuthProcessorStartOAuthSignInResult, error) {
	startResult, err := self.startSignIn(ctx, params.AuthProcessorStartSignInParams)
	if err != nil {
		return nil, err
	}

	oAuthProvider, err := goth.GetProvider(strings.ToLower(params.Provider))
	if err != nil {
		return nil, ErrAuthProcessorInvalidOAuthProvider.Raise().Cause(err)
	}

	oAuthSession, err := oAuthProvider.BeginAuth(url.QueryEscape(startResult.State))
	if err != nil {
		return nil, ErrAuthProcessorGeneric.Raise().Cause(err)
	}

	authURL, err := oAuthSession.GetAuthURL()
	if err != nil {
		return nil, ErrAuthProcessorGeneric.Raise().Cause(err)
	}

	result := AuthProcessorStartOAuthSignInResult{}
	result.AuthProcessorStartSignInResult = *startResult
	result.AuthURL = authURL

	return &result, nil
}

type AuthProcessorStartSAMLSignInParams struct {
	AuthProcessorStartSignInParams
}

type AuthProcessorStartSAMLSignInResult struct {
	AuthProcessorStartSignInResult
}

func (self *AuthProcessor) StartSAMLSignIn(ctx context.Context,
	params AuthProcessorStartSAMLSignInParams) (*AuthProcessorStartSAMLSignInResult, error) {
	return nil, ErrAuthProcessorGeneric.Raise().With("SAML not implemented")
}

type AuthProcessorEndSignInParams struct {
	Provider   string
	StartState string
	EndState   string
	IP         string
	Device     string
	Name       string
	Picture    string
	Email      string
}

type AuthProcessorEndSignInResult struct {
	Token      string
	RedirectTo *string
}

func (self *AuthProcessor) endSignIn(ctx context.Context,
	params AuthProcessorEndSignInParams) (*AuthProcessorEndSignInResult, error) {
	var state State

	if params.StartState != params.EndState {
		return nil, ErrAuthProcessorInvalidState.Raise().With("state doesn't match")
	}

	err := self.tokenizer.Decrypt(params.StartState, pvx.NewSymmetricKey(self.cryptKey, pvx.Version4)).Scan(&state, nil)
	if err != nil {
		return nil, ErrAuthProcessorInvalidState.Raise().Cause(err)
	}

	if time.Now().After(state.ExpiresAt) {
		return nil, ErrAuthProcessorInvalidState.Raise().With("state expired")
	}

	if state.RedirectTo != nil && !util.SameOrigin(*state.RedirectTo, self.config.Frontend.BaseURL) {
		return nil, ErrAuthProcessorInvalidState.Raise().With("can only redirect within frontend")
	}

	var token string

	err = self.database.Transaction(ctx, nil, func(ctx context.Context) error {
		var session *Session
		var _user *user.User

		if state.PreviousToken != nil {
			verifierResult, err := self.authVerifier.CheckToken(ctx, AuthVerifierCheckTokenParams{
				Token: *state.PreviousToken,
			})
			if err != nil && !ErrAuthVerifierInvalidToken.Is(err) {
				if ErrAuthVerifierUnauthorizedUser.Is(err) {
					return ErrAuthProcessorUnauthorizedUser.Raise().Cause(err)
				}

				return ErrAuthProcessorInvalidState.Raise().Cause(err)
			}

			if err == nil {
				if params.Email != verifierResult.User.Email {
					// TODO: In the future allow to add the new email to the existing user
					return ErrAuthProcessorInvalidState.Raise().
						With("multiple emails for the same user not supported")
				}

				if params.Provider != verifierResult.Session.Provider {
					verifierResult.Session.Provider = params.Provider

					err := self.sessionRepository.UpdateProvider(ctx, verifierResult.Session.ID,
						verifierResult.Session.Provider)
					if err != nil {
						return ErrAuthProcessorInvalidState.Raise().Cause(err)
					}
				}

				session = &verifierResult.Session
				_user = &verifierResult.User
			}
		}

		if _user == nil {
			_user, err = self.userRepository.GetByEmail(ctx, params.Email)
			if err != nil {
				return ErrAuthProcessorGeneric.Raise().Cause(err)
			}

			if _user != nil && _user.DeletedAt != nil {
				err := self.userRepository.UpdateEmail(ctx, _user.ID, _user.ID+"@deleted.user")
				if err != nil {
					return ErrAuthProcessorGeneric.Raise().Cause(err)
				}

				_user = nil
			}

			if _user == nil {
				_user = user.NewUser()
				_user.ID = xid.New().String()
				_user.OrganizationID = ""
				_user.Name = params.Name
				_user.Picture = params.Picture
				_user.Email = params.Email
				_user.Role = user.UserRoleMember
				_user.Settings = user.UserSettings{}
				_user.CreatedAt = time.Now()
				_user.DeletedAt = nil

				invitation, err := self.invitationRepository.GetByEmail(ctx, _user.Email)
				if err != nil {
					return ErrAuthProcessorGeneric.Raise().Cause(err)
				}

				if invitation != nil {
					if time.Now().Before(invitation.ExpiresAt) {
						_organization, err := self.organizationRepository.GetByID(ctx, invitation.OrganizationID)
						if err != nil {
							return ErrAuthProcessorGeneric.Raise().Cause(err)
						}

						if _organization == nil {
							return ErrAuthProcessorGeneric.Raise().With("no organization linked")
						}

						if _organization.DeletedAt != nil {
							return ErrAuthProcessorUnauthorizedUser.Raise().With("organization is deleted")
						}

						_user.OrganizationID = invitation.OrganizationID
						_user.Role = invitation.Role
					}

					err = self.invitationRepository.DeleteByID(ctx, invitation.ID)
					if err != nil {
						return ErrAuthProcessorGeneric.Raise().Cause(err)
					}
				}

				if len(_user.OrganizationID) == 0 {
					domain := strings.Split(_user.Email, "@")[1]

					_organization, err := self.organizationRepository.GetByDomain(ctx, domain)
					if err != nil {
						return ErrAuthProcessorGeneric.Raise().Cause(err)
					}

					if _organization != nil {
						if _organization.DeletedAt != nil {
							err := self.organizationRepository.UpdateDomain(ctx, _organization.ID,
								_organization.ID+".deleted.organization")
							if err != nil {
								return ErrAuthProcessorGeneric.Raise().Cause(err)
							}
						} else if _organization.Settings.DomainSignIn &&
							organization.IsDomainSignInSupported(_organization.Domain) {
							_user.OrganizationID = _organization.ID
							_user.Role = user.UserRoleMember
						}
					}
				}

				if len(_user.OrganizationID) == 0 {
					domain := strings.Split(_user.Email, "@")[1]
					domainSignIn := true
					if !organization.IsDomainSignInSupported(domain) {
						domain = strings.ReplaceAll(_user.Email, "@", ".")
						domainSignIn = false
					}
					domainParts := strings.Split(domain, ".")
					// TODO: Fix this doesn't do anything as we have lost the original domain info!
					included, extra := organization.PlanTrialCapacity(domain)
					now := time.Now()

					_organization := organization.NewOrganization()
					_organization.ID = xid.New().String()
					_organization.Name = strings.Join(domainParts[:len(domainParts)-1], ".")
					_organization.Picture = organization.ORGANIZATION_DEFAULT_PICTURE
					_organization.Domain = domain
					_organization.Settings.DomainSignIn = domainSignIn
					_organization.Plan = organization.OrganizationPlanTrial
					_organization.TrialEndsAt = now.Add(organization.ORGANIZATION_PLAN_TRIAL_PERIOD)
					_organization.Capacity.Included = included
					_organization.Capacity.Extra = extra
					_organization.Usage.Value = 0
					_organization.Usage.LastComputedAt = now
					_organization.CreatedAt = now
					_organization.DeletedAt = nil

					_organization, err := self.organizationRepository.Create(ctx, *_organization)
					if err != nil {
						return ErrAuthProcessorGeneric.Raise().Cause(err)
					}

					_user.OrganizationID = _organization.ID
					_user.Role = user.UserRoleAdmin
				}

				_user, err = self.userRepository.Create(ctx, *_user)
				if err != nil {
					return ErrAuthProcessorGeneric.Raise().Cause(err)
				}
			}

			session = NewSession()
			session.ID = xid.New().String()
			session.UserID = _user.ID
			session.Provider = params.Provider
			session.Metadata.Locations = []SessionMetadataLocation{{IP: params.IP, Device: params.Device}}
			session.CreatedAt = time.Now()
			session.LastSeenAt = session.CreatedAt
			session.ExpiredAt = nil

			session, err = self.sessionRepository.Create(ctx, *session)
			if err != nil {
				return ErrAuthProcessorGeneric.Raise().Cause(err)
			}
		}

		tokenPayload := NewToken()
		tokenPayload.Private.SessionID = session.ID
		tokenPayload.Private.Provider = session.Provider
		tokenPayload.Private.ExpiresAt = time.Now().Add(TOKEN_EXPIRATION)
		tokenPayload.Private.Nonce = kitUtil.RandomString(8)

		token, err = self.tokenizer.Encrypt(pvx.NewSymmetricKey(self.cryptKey, pvx.Version4),
			&tokenPayload.Private, pvx.WithFooter(&tokenPayload.Public))
		if err != nil {
			return ErrAuthProcessorGeneric.Raise().Cause(err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	result := AuthProcessorEndSignInResult{}
	result.Token = token
	result.RedirectTo = state.RedirectTo

	return &result, nil
}

type AuthProcessorEndEmailSignInParams struct {
	AuthProcessorEndSignInParams
	SignInCodeID   string
	SignInCodeCode string
}

type AuthProcessorEndEmailSignInResult struct {
	AuthProcessorEndSignInResult
}

func (self *AuthProcessor) EndEmailSignIn(ctx context.Context,
	params AuthProcessorEndEmailSignInParams) (*AuthProcessorEndEmailSignInResult, error) {
	signInCode, err := self.signInCodeRepository.GetByID(ctx, params.SignInCodeID)
	if err != nil {
		return nil, ErrAuthProcessorGeneric.Raise().Cause(err)
	}

	if signInCode == nil {
		return nil, ErrAuthProcessorInvalidSignInCode.Raise().With("no sign in code linked")
	}

	if time.Now().After(signInCode.ExpiresAt) {
		return nil, ErrAuthProcessorInvalidSignInCode.Raise().With("expired")
	}

	signInCode.Attempts++

	err = self.signInCodeRepository.UpdateAttempts(ctx, signInCode.ID, signInCode.Attempts)
	if err != nil {
		return nil, ErrAuthProcessorGeneric.Raise().Cause(err)
	}

	if signInCode.Attempts > SIGN_IN_CODE_MAX_ATTEMPTS {
		return nil, ErrAuthProcessorInvalidSignInCode.Raise().With("max attempts reached")
	}

	if params.SignInCodeCode != signInCode.Code {
		return nil, ErrAuthProcessorInvalidSignInCode.Raise().With("code doesn't match")
	}

	err = self.signInCodeRepository.DeleteByID(ctx, signInCode.ID)
	if err != nil {
		return nil, ErrAuthProcessorGeneric.Raise().Cause(err)
	}

	params.AuthProcessorEndSignInParams.Provider = SessionProviderEmail
	params.AuthProcessorEndSignInParams.Name = strings.Split(signInCode.Email, "@")[0]
	params.AuthProcessorEndSignInParams.Picture = user.USER_DEFAULT_PICTURE
	params.AuthProcessorEndSignInParams.Email = signInCode.Email

	endResult, err := self.endSignIn(ctx, params.AuthProcessorEndSignInParams)
	if err != nil {
		return nil, err
	}

	result := AuthProcessorEndEmailSignInResult{}
	result.AuthProcessorEndSignInResult = *endResult

	return &result, nil
}

type AuthProcessorEndOAuthSignInParams struct {
	AuthProcessorEndSignInParams
	Provider   string
	AuthResult string
}

type AuthProcessorEndOAuthSignInResult struct {
	AuthProcessorEndSignInResult
}

func (self *AuthProcessor) EndOAuthSignIn(ctx context.Context,
	params AuthProcessorEndOAuthSignInParams) (*AuthProcessorEndOAuthSignInResult, error) {
	oAuthProvider, err := goth.GetProvider(strings.ToLower(params.Provider))
	if err != nil {
		return nil, ErrAuthProcessorInvalidOAuthProvider.Raise().Cause(err)
	}

	oAuthSession, err := oAuthProvider.BeginAuth(url.QueryEscape(params.StartState))
	if err != nil {
		return nil, ErrAuthProcessorGeneric.Raise().Cause(err)
	}

	oAuthValues, err := url.ParseQuery(params.AuthResult)
	if err != nil {
		return nil, ErrAuthProcessorInvalidOAuthResult.Raise().Cause(err)
	}

	_, err = oAuthSession.Authorize(oAuthProvider, oAuthValues)
	if err != nil {
		return nil, ErrAuthProcessorInvalidOAuthResult.Raise().Cause(err)
	}

	oAuthUser, err := oAuthProvider.FetchUser(oAuthSession)
	if err != nil {
		return nil, ErrAuthProcessorGeneric.Raise().Cause(err)
	}

	name := oAuthUser.Name
	if len(name) == 0 {
		name = oAuthUser.FirstName
		if len(name) == 0 {
			name = oAuthUser.NickName
			if len(name) == 0 {
				name = strings.Split(oAuthUser.Email, "@")[0]
			}
		}
	}

	picture := oAuthUser.AvatarURL
	if len(picture) == 0 {
		picture = user.USER_DEFAULT_PICTURE
	}

	params.AuthProcessorEndSignInParams.Provider = params.Provider
	params.AuthProcessorEndSignInParams.Name = name
	params.AuthProcessorEndSignInParams.Picture = picture
	params.AuthProcessorEndSignInParams.Email = oAuthUser.Email

	endResult, err := self.endSignIn(ctx, params.AuthProcessorEndSignInParams)
	if err != nil {
		return nil, err
	}

	result := AuthProcessorEndOAuthSignInResult{}
	result.AuthProcessorEndSignInResult = *endResult

	return &result, nil
}

type AuthProcessorEndSAMLSignInParams struct {
	AuthProcessorEndSignInParams
}

type AuthProcessorEndSAMLSignInResult struct {
	AuthProcessorEndSignInResult
}

func (self *AuthProcessor) EndSAMLSignIn(ctx context.Context,
	params AuthProcessorEndSAMLSignInParams) (*AuthProcessorEndSAMLSignInResult, error) {
	return nil, ErrAuthProcessorGeneric.Raise().With("SAML not implemented")
}

type AuthProcessorSignOutParams struct {
	Session Session
}

func (self *AuthProcessor) SignOut(ctx context.Context, params AuthProcessorSignOutParams) error {
	err := self.sessionRepository.UpdateExpiredAt(ctx, params.Session.ID, time.Now())
	if err != nil {
		return ErrAuthProcessorGeneric.Raise().Cause(err)
	}

	return nil
}
