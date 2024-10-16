package auth

import (
	"context"
	"crypto/sha256"
	"time"

	"backend/pkg/config"
	"backend/pkg/organization"
	"backend/pkg/user"

	"github.com/neoxelox/errors"
	"github.com/neoxelox/kit"
	"github.com/vk-rv/pvx"
)

var (
	ErrAuthVerifierGeneric          = errors.New("auth verifier failed")
	ErrAuthVerifierInvalidToken     = errors.New("token is invalid")
	ErrAuthVerifierUnauthorizedUser = errors.New("user is unauthorized")
)

type AuthVerifier struct {
	config                 config.Config
	observer               *kit.Observer
	cryptKey               []byte
	tokenizer              *pvx.ProtoV4Local
	sessionRepository      SessionRepository
	userRepository         user.UserRepository
	organizationRepository organization.OrganizationRepository
}

func NewAuthVerifier(observer *kit.Observer, sessionRepository SessionRepository,
	userRepository user.UserRepository, organizationRepository organization.OrganizationRepository,
	config config.Config) *AuthVerifier {
	cryptKey := sha256.Sum256([]byte(config.Auth.CryptKey))
	tokenizer := pvx.NewPV4Local()

	return &AuthVerifier{
		config:                 config,
		observer:               observer,
		cryptKey:               cryptKey[:],
		tokenizer:              tokenizer,
		sessionRepository:      sessionRepository,
		userRepository:         userRepository,
		organizationRepository: organizationRepository,
	}
}

type AuthVerifierCheckTokenParams struct {
	Token string
}

type AuthVerifierCheckTokenResult struct {
	Session      Session
	User         user.User
	Organization organization.Organization
}

func (self *AuthVerifier) CheckToken(ctx context.Context,
	params AuthVerifierCheckTokenParams) (*AuthVerifierCheckTokenResult, error) {
	var token Token

	err := self.tokenizer.Decrypt(params.Token, pvx.NewSymmetricKey(self.cryptKey, pvx.Version4)).
		Scan(&token.Private, &token.Public)
	if err != nil {
		return nil, ErrAuthVerifierInvalidToken.Raise().Cause(err)
	}

	if time.Now().After(token.Private.ExpiresAt) {
		return nil, ErrAuthVerifierInvalidToken.Raise().With("token expired")
	}

	session, err := self.sessionRepository.GetByID(ctx, token.Private.SessionID)
	if err != nil {
		return nil, ErrAuthVerifierGeneric.Raise().Cause(err)
	}

	if session == nil {
		return nil, ErrAuthVerifierInvalidToken.Raise().With("no session linked")
	}

	if session.ExpiredAt != nil && time.Now().After(*session.ExpiredAt) {
		return nil, ErrAuthVerifierInvalidToken.Raise().With("session expired")
	}

	if token.Private.Provider != session.Provider {
		return nil, ErrAuthVerifierInvalidToken.Raise().With("provider doesn't match")
	}

	user, err := self.userRepository.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, ErrAuthVerifierGeneric.Raise().Cause(err)
	}

	if user == nil {
		return nil, ErrAuthVerifierInvalidToken.Raise().With("no user linked")
	}

	if user.DeletedAt != nil {
		return nil, ErrAuthVerifierUnauthorizedUser.Raise().With("user is deleted")
	}

	organization, err := self.organizationRepository.GetByID(ctx, user.OrganizationID)
	if err != nil {
		return nil, ErrAuthVerifierGeneric.Raise().Cause(err)
	}

	if organization == nil {
		return nil, ErrAuthVerifierInvalidToken.Raise().With("no organization linked")
	}

	if organization.DeletedAt != nil {
		return nil, ErrAuthVerifierUnauthorizedUser.Raise().With("organization is deleted")
	}

	result := AuthVerifierCheckTokenResult{}
	result.Session = *session
	result.User = *user
	result.Organization = *organization

	return &result, nil
}
