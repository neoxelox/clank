package auth

import (
	"fmt"
	"time"

	"github.com/neoxelox/kit/util"
)

const (
	SessionProviderEmail  = "EMAIL"
	SessionProviderGoogle = "GOOGLE"
	SessionProviderApple  = "APPLE"
	SessionProviderAmazon = "AMAZON"
	SessionProviderSAML   = "SAML"
)

func IsSessionProvider(value string) bool {
	return value == SessionProviderEmail ||
		value == SessionProviderGoogle ||
		value == SessionProviderApple ||
		value == SessionProviderAmazon ||
		value == SessionProviderSAML
}

type SessionMetadataLocation struct {
	IP     string
	Device string
}

type SessionMetadata struct {
	Locations []SessionMetadataLocation
}

type Session struct {
	ID         string
	UserID     string
	Provider   string
	Metadata   SessionMetadata
	CreatedAt  time.Time
	LastSeenAt time.Time
	ExpiredAt  *time.Time
}

func NewSession() *Session {
	return &Session{}
}

func (self Session) String() string {
	return fmt.Sprintf("<Session: %s (%s)>", self.UserID, self.ID)
}

func (self Session) Equals(other Session) bool {
	return util.Equals(self, other)
}

func (self Session) Copy() *Session {
	return util.Copy(self)
}

const (
	STATE_COOKIE_NAME = "_clank_state"
	STATE_EXPIRATION  = 5 * time.Minute
)

type State struct {
	RedirectTo    *string
	PreviousToken *string
	ExpiresAt     time.Time
	Nonce         string
}

func NewState() *State {
	return &State{}
}

// Needed for Paseto library (Validation will be performed in the usecase)
func (self *State) Valid() error { return nil }

func (self State) String() string {
	return fmt.Sprintf("<State: %s (%s)>", self.Nonce, self.ExpiresAt)
}

func (self State) Equals(other State) bool {
	return util.Equals(self, other)
}

func (self State) Copy() *State {
	return util.Copy(self)
}

const (
	TOKEN_COOKIE_NAME = "_clank_token"
	TOKEN_EXPIRATION  = 30 * 24 * time.Hour
)

type TokenPrivate struct {
	SessionID string
	Provider  string
	ExpiresAt time.Time
	Nonce     string
}

// Needed for Paseto library (Validation will be performed in the usecase)
func (self *TokenPrivate) Valid() error { return nil }

type TokenPublic struct {
}

type Token struct {
	Private TokenPrivate
	Public  TokenPublic
}

func NewToken() *Token {
	return &Token{}
}

func (self Token) String() string {
	return fmt.Sprintf("<Token: %s (%s)>", self.Private.SessionID, self.Private.ExpiresAt)
}

func (self Token) Equals(other Token) bool {
	return util.Equals(self, other)
}

func (self Token) Copy() *Token {
	return util.Copy(self)
}

const (
	SIGN_IN_CODE_LENGTH       = 6
	SIGN_IN_CODE_MAX_ATTEMPTS = 3
	SIGN_IN_CODE_EXPIRATION   = 5 * time.Minute
)

type SignInCode struct {
	ID        string
	Email     string
	Code      string
	Attempts  int
	ExpiresAt time.Time
}

func NewSignInCode() *SignInCode {
	return &SignInCode{}
}

func (self SignInCode) String() string {
	return fmt.Sprintf("<SignInCode: %s (%s)>", self.Email, self.ID)
}

func (self SignInCode) Equals(other SignInCode) bool {
	return util.Equals(self, other)
}

func (self SignInCode) Copy() *SignInCode {
	return util.Copy(self)
}
