package auth_test

import (
	"context"
	"crypto/sha256"
	"strings"
	"testing"
	"time"

	"backend/pkg/auth"
	"backend/pkg/config"
	"backend/pkg/organization"
	"backend/pkg/user"

	"github.com/neoxelox/kit"
	"github.com/neoxelox/kit/util"
	"github.com/stretchr/testify/suite"
	"github.com/vk-rv/pvx"
)

type AuthVerifierTestSuite struct {
	suite.Suite
	ctx    context.Context
	config config.Config
	mocks  struct {
		sessionRepository      *auth.SessionRepositoryMock
		userRepository         *user.UserRepositoryMock
		organizationRepository *organization.OrganizationRepositoryMock
	}
	verifier *auth.AuthVerifier
}

func (self *AuthVerifierTestSuite) SetupTest() {
	self.ctx = context.Background()

	self.config = *config.NewConfig()
	self.config.Service.Environment = kit.EnvIntegration
	self.config.Service.Release = "test"
	self.config.Service.Name = "test"
	self.config.Auth.CryptKey = "test-crypt-key"

	observer, err := kit.NewObserver(self.ctx, kit.ObserverConfig{
		Environment: self.config.Service.Environment,
		Release:     self.config.Service.Release,
		Service:     self.config.Service.Name,
		Level:       kit.LvlError,
	})
	self.Require().NoError(err)

	self.mocks.sessionRepository = auth.NewSessionRepositoryMock()
	self.mocks.userRepository = user.NewUserRepositoryMock()
	self.mocks.organizationRepository = organization.NewOrganizationRepositoryMock()

	self.verifier = auth.NewAuthVerifier(observer, self.mocks.sessionRepository,
		self.mocks.userRepository, self.mocks.organizationRepository, self.config)
}

func (self *AuthVerifierTestSuite) createToken(payload *auth.Token) string {
	payload.Private.Nonce = "test-nonce"

	cryptKey := sha256.Sum256([]byte(self.config.Auth.CryptKey))
	token, err := pvx.NewPV4Local().Encrypt(
		pvx.NewSymmetricKey(cryptKey[:], pvx.Version4), &payload.Private, pvx.WithFooter(&payload.Public))
	self.Require().NoError(err)

	return token
}

func TestAuthVerifierSuite(t *testing.T) {
	suite.Run(t, new(AuthVerifierTestSuite))
}

func (self *AuthVerifierTestSuite) TestValidToken() {
	// Given: A valid token with a valid linked session, user and organization
	self.mocks.sessionRepository.On("GetByID", self.ctx, "session1").Return(&auth.Session{
		ID:        "session1",
		UserID:    "user1",
		Provider:  "test-provider",
		ExpiredAt: nil,
	}, nil)
	self.mocks.userRepository.On("GetByID", self.ctx, "user1").Return(&user.User{
		ID:             "user1",
		OrganizationID: "org1",
		DeletedAt:      nil,
	}, nil)
	self.mocks.organizationRepository.On("GetByID", self.ctx, "org1").Return(&organization.Organization{
		ID:        "org1",
		DeletedAt: nil,
	}, nil)
	token := self.createToken(&auth.Token{
		Private: auth.TokenPrivate{
			SessionID: "session1",
			Provider:  "test-provider",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		},
	})

	// When: Checking the token
	result, err := self.verifier.CheckToken(self.ctx, auth.AuthVerifierCheckTokenParams{Token: token})

	// Then: The token is valid and the linked data is returned
	self.Require().NoError(err)
	self.Require().Equal("session1", result.Session.ID)
	self.Require().Equal("user1", result.User.ID)
	self.Require().Equal("org1", result.Organization.ID)
}

func (self *AuthVerifierTestSuite) TestInvalidToken() {
	// Given: An invalid token
	token := "invalid-token"

	// When: Checking the token
	result, err := self.verifier.CheckToken(self.ctx, auth.AuthVerifierCheckTokenParams{Token: token})

	// Then: The token is invalid
	self.Require().True(auth.ErrAuthVerifierInvalidToken.Is(err))
	self.Require().Nil(result)
}

func (self *AuthVerifierTestSuite) TestExpiredToken() {
	// Given: An expired token
	expiredToken := self.createToken(&auth.Token{
		Private: auth.TokenPrivate{
			SessionID: "session1",
			Provider:  "test-provider",
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		},
	})

	// When: Checking the token
	result, err := self.verifier.CheckToken(self.ctx, auth.AuthVerifierCheckTokenParams{Token: expiredToken})

	// Then: The token is invalid
	self.Require().True(auth.ErrAuthVerifierInvalidToken.Is(err))
	self.Require().True(strings.Contains(err.Error(), "token expired"))
	self.Require().Nil(result)
}

func (self *AuthVerifierTestSuite) TestSessionNotFound() {
	// Given: A token with a linked session that does not exist
	self.mocks.sessionRepository.On("GetByID", self.ctx, "session1").Return(nil, nil)
	token := self.createToken(&auth.Token{
		Private: auth.TokenPrivate{
			SessionID: "session1",
			Provider:  "test-provider",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		},
	})

	// When: Checking the token
	result, err := self.verifier.CheckToken(self.ctx, auth.AuthVerifierCheckTokenParams{Token: token})

	// Then: The token is invalid
	self.Require().True(auth.ErrAuthVerifierInvalidToken.Is(err))
	self.Require().True(strings.Contains(err.Error(), "no session linked"))
	self.Require().Nil(result)
}

func (self *AuthVerifierTestSuite) TestExpiredSession() {
	// Given: A token with a linked session that has expired
	self.mocks.sessionRepository.On("GetByID", self.ctx, "session1").Return(&auth.Session{
		ID:        "session1",
		UserID:    "user1",
		Provider:  "test-provider",
		ExpiredAt: util.Pointer(time.Now().Add(-1 * time.Hour)),
	}, nil)
	token := self.createToken(&auth.Token{
		Private: auth.TokenPrivate{
			SessionID: "session1",
			Provider:  "test-provider",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		},
	})

	// When: Checking the token
	result, err := self.verifier.CheckToken(self.ctx, auth.AuthVerifierCheckTokenParams{Token: token})

	// Then: The token is invalid
	self.Require().True(auth.ErrAuthVerifierInvalidToken.Is(err))
	self.Require().True(strings.Contains(err.Error(), "session expired"))
	self.Require().Nil(result)
}

func (self *AuthVerifierTestSuite) TestProviderMismatch() {
	// Given: A token with a linked session that has a different provider than the token's provider
	self.mocks.sessionRepository.On("GetByID", self.ctx, "session1").Return(&auth.Session{
		ID:       "session1",
		UserID:   "user1",
		Provider: "different-provider",
	}, nil)
	token := self.createToken(&auth.Token{
		Private: auth.TokenPrivate{
			SessionID: "session1",
			Provider:  "test-provider",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		},
	})

	// When: Checking the token
	result, err := self.verifier.CheckToken(self.ctx, auth.AuthVerifierCheckTokenParams{Token: token})

	// Then: The token is invalid
	self.Require().True(auth.ErrAuthVerifierInvalidToken.Is(err))
	self.Require().True(strings.Contains(err.Error(), "provider doesn't match"))
	self.Require().Nil(result)
}

func (self *AuthVerifierTestSuite) TestUserNotFound() {
	// Given: A token with a linked user that does not exist
	self.mocks.sessionRepository.On("GetByID", self.ctx, "session1").Return(&auth.Session{
		ID:       "session1",
		UserID:   "user1",
		Provider: "test-provider",
	}, nil)
	self.mocks.userRepository.On("GetByID", self.ctx, "user1").Return(nil, nil)
	token := self.createToken(&auth.Token{
		Private: auth.TokenPrivate{
			SessionID: "session1",
			Provider:  "test-provider",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		},
	})

	// When: Checking the token
	result, err := self.verifier.CheckToken(self.ctx, auth.AuthVerifierCheckTokenParams{Token: token})

	// Then: The token is invalid
	self.Require().True(auth.ErrAuthVerifierInvalidToken.Is(err))
	self.Require().True(strings.Contains(err.Error(), "no user linked"))
	self.Require().Nil(result)
}

func (self *AuthVerifierTestSuite) TestDeletedUser() {
	// Given: A token with a linked user that has been deleted
	self.mocks.sessionRepository.On("GetByID", self.ctx, "session1").Return(&auth.Session{
		ID:       "session1",
		UserID:   "user1",
		Provider: "test-provider",
	}, nil)
	self.mocks.userRepository.On("GetByID", self.ctx, "user1").Return(&user.User{
		ID:             "user1",
		OrganizationID: "org1",
		DeletedAt:      util.Pointer(time.Now().Add(-1 * time.Hour)),
	}, nil)
	token := self.createToken(&auth.Token{
		Private: auth.TokenPrivate{
			SessionID: "session1",
			Provider:  "test-provider",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		},
	})

	// When: Checking the token
	result, err := self.verifier.CheckToken(self.ctx, auth.AuthVerifierCheckTokenParams{Token: token})

	// Then: The token is invalid
	self.Require().True(auth.ErrAuthVerifierUnauthorizedUser.Is(err))
	self.Require().True(strings.Contains(err.Error(), "user is deleted"))
	self.Require().Nil(result)
}

func (self *AuthVerifierTestSuite) TestOrganizationNotFound() {
	// Given: A token with a linked organization that does not exist
	self.mocks.sessionRepository.On("GetByID", self.ctx, "session1").Return(&auth.Session{
		ID:       "session1",
		UserID:   "user1",
		Provider: "test-provider",
	}, nil)
	self.mocks.userRepository.On("GetByID", self.ctx, "user1").Return(&user.User{
		ID:             "user1",
		OrganizationID: "org1",
	}, nil)
	self.mocks.organizationRepository.On("GetByID", self.ctx, "org1").Return(nil, nil)
	token := self.createToken(&auth.Token{
		Private: auth.TokenPrivate{
			SessionID: "session1",
			Provider:  "test-provider",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		},
	})

	// When: Checking the token
	result, err := self.verifier.CheckToken(self.ctx, auth.AuthVerifierCheckTokenParams{Token: token})

	// Then: The token is invalid
	self.Require().True(auth.ErrAuthVerifierInvalidToken.Is(err))
	self.Require().True(strings.Contains(err.Error(), "no organization linked"))
	self.Require().Nil(result)
}

func (self *AuthVerifierTestSuite) TestDeletedOrganization() {
	// Given: A token with a linked organization that has been deleted
	self.mocks.sessionRepository.On("GetByID", self.ctx, "session1").Return(&auth.Session{
		ID:       "session1",
		UserID:   "user1",
		Provider: "test-provider",
	}, nil)
	self.mocks.userRepository.On("GetByID", self.ctx, "user1").Return(&user.User{
		ID:             "user1",
		OrganizationID: "org1",
	}, nil)
	self.mocks.organizationRepository.On("GetByID", self.ctx, "org1").Return(&organization.Organization{
		ID:        "org1",
		DeletedAt: util.Pointer(time.Now().Add(-1 * time.Hour)),
	}, nil)
	token := self.createToken(&auth.Token{
		Private: auth.TokenPrivate{
			SessionID: "session1",
			Provider:  "test-provider",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		},
	})

	// When: Checking the token
	result, err := self.verifier.CheckToken(self.ctx, auth.AuthVerifierCheckTokenParams{Token: token})

	// Then: The token is invalid
	self.Require().True(auth.ErrAuthVerifierUnauthorizedUser.Is(err))
	self.Require().True(strings.Contains(err.Error(), "organization is deleted"))
	self.Require().Nil(result)
}
