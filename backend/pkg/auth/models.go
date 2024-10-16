package auth

import (
	"encoding/json"
	"time"
)

const (
	SESSION_MODEL_TABLE = "\"session\""
)

type SessionModel struct {
	ID         string     `db:"id"`
	UserID     string     `db:"user_id"`
	Provider   string     `db:"provider"`
	Metadata   []byte     `db:"metadata"`
	CreatedAt  time.Time  `db:"created_at"`
	LastSeenAt time.Time  `db:"last_seen_at"`
	ExpiredAt  *time.Time `db:"expired_at"`
}

func NewSessionModel(session Session) *SessionModel {
	metadata, err := json.Marshal(session.Metadata)
	if err != nil {
		panic(err)
	}

	return &SessionModel{
		ID:         session.ID,
		UserID:     session.UserID,
		Provider:   session.Provider,
		Metadata:   metadata,
		CreatedAt:  session.CreatedAt,
		LastSeenAt: session.LastSeenAt,
		ExpiredAt:  session.ExpiredAt,
	}
}

func (self *SessionModel) ToEntity() *Session {
	var metadata SessionMetadata
	err := json.Unmarshal(self.Metadata, &metadata)
	if err != nil {
		panic(err)
	}

	return &Session{
		ID:         self.ID,
		UserID:     self.UserID,
		Provider:   self.Provider,
		Metadata:   metadata,
		CreatedAt:  self.CreatedAt,
		LastSeenAt: self.LastSeenAt,
		ExpiredAt:  self.ExpiredAt,
	}
}

const (
	SIGN_IN_CODE_MODEL_TABLE = "\"sign_in_code\""
)

type SignInCodeModel struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Code      string    `db:"code"`
	Attempts  int       `db:"attempts"`
	ExpiresAt time.Time `db:"expires_at"`
}

func NewSignInCodeModel(signInCode SignInCode) *SignInCodeModel {
	return &SignInCodeModel{
		ID:        signInCode.ID,
		Email:     signInCode.Email,
		Code:      signInCode.Code,
		Attempts:  signInCode.Attempts,
		ExpiresAt: signInCode.ExpiresAt,
	}
}

func (self *SignInCodeModel) ToEntity() *SignInCode {
	return &SignInCode{
		ID:        self.ID,
		Email:     self.Email,
		Code:      self.Code,
		Attempts:  self.Attempts,
		ExpiresAt: self.ExpiresAt,
	}
}
