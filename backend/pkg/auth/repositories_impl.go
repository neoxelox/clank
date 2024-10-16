package auth

import (
	"context"
	"time"

	"github.com/leporo/sqlf"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
)

type SessionRepositoryImpl struct {
	config   config.Config
	observer *kit.Observer
	database *kit.Database
}

func NewSessionRepositoryImpl(observer *kit.Observer, database *kit.Database, config config.Config) *SessionRepositoryImpl {
	return &SessionRepositoryImpl{
		config:   config,
		observer: observer,
		database: database,
	}
}

func (self *SessionRepositoryImpl) Create(ctx context.Context, session Session) (*Session, error) {
	s := NewSessionModel(session)

	stmt := sqlf.
		InsertInto(SESSION_MODEL_TABLE).
		Set("id", s.ID).
		Set("user_id", s.UserID).
		Set("provider", s.Provider).
		Set("metadata", s.Metadata).
		Set("created_at", s.CreatedAt).
		Set("last_seen_at", s.LastSeenAt).
		Set("expired_at", s.ExpiredAt).
		Returning("*").To(&s)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	return s.ToEntity(), nil
}

func (self *SessionRepositoryImpl) GetByID(ctx context.Context, id string) (*Session, error) {
	var s SessionModel

	stmt := sqlf.
		Select("*").To(&s).
		From(SESSION_MODEL_TABLE).
		Where("id = ?", id)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return s.ToEntity(), nil
}

func (self *SessionRepositoryImpl) UpdateProvider(ctx context.Context, id string, provider string) error {
	stmt := sqlf.
		Update(SESSION_MODEL_TABLE).
		Set("provider", provider).
		Where("id = ?", id)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

func (self *SessionRepositoryImpl) UpdateSeen(ctx context.Context, session Session) error {
	s := NewSessionModel(session)

	stmt := sqlf.
		Update(SESSION_MODEL_TABLE).
		Set("metadata", s.Metadata).
		Set("last_seen_at", s.LastSeenAt).
		Where("id = ?", s.ID)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

func (self *SessionRepositoryImpl) UpdateExpiredAt(ctx context.Context, id string, expiredAt time.Time) error {
	stmt := sqlf.
		Update(SESSION_MODEL_TABLE).
		Set("expired_at", expiredAt).
		Where("id = ?", id)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

type SignInCodeRepositoryImpl struct {
	config   config.Config
	observer *kit.Observer
	database *kit.Database
}

func NewSignInCodeRepositoryImpl(observer *kit.Observer, database *kit.Database,
	config config.Config) *SignInCodeRepositoryImpl {
	return &SignInCodeRepositoryImpl{
		config:   config,
		observer: observer,
		database: database,
	}
}

func (self *SignInCodeRepositoryImpl) Create(ctx context.Context, signInCode SignInCode) (*SignInCode, error) {
	s := NewSignInCodeModel(signInCode)

	stmt := sqlf.
		InsertInto(SIGN_IN_CODE_MODEL_TABLE).
		Set("id", s.ID).
		Set("email", s.Email).
		Set("code", s.Code).
		Set("attempts", s.Attempts).
		Set("expires_at", s.ExpiresAt).
		Returning("*").To(&s)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	return s.ToEntity(), nil
}

func (self *SignInCodeRepositoryImpl) GetByID(ctx context.Context, id string) (*SignInCode, error) {
	var s SignInCodeModel

	stmt := sqlf.
		Select("*").To(&s).
		From(SIGN_IN_CODE_MODEL_TABLE).
		Where("id = ?", id)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return s.ToEntity(), nil
}

func (self *SignInCodeRepositoryImpl) GetByEmail(ctx context.Context, email string) (*SignInCode, error) {
	var s SignInCodeModel

	stmt := sqlf.
		Select("*").To(&s).
		From(SIGN_IN_CODE_MODEL_TABLE).
		Where("email = ?", email)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return s.ToEntity(), nil
}

func (self *SignInCodeRepositoryImpl) UpdateAttempts(ctx context.Context, id string, attempts int) error {
	stmt := sqlf.
		Update(SIGN_IN_CODE_MODEL_TABLE).
		Set("attempts", attempts).
		Where("id = ?", id)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

func (self *SignInCodeRepositoryImpl) DeleteByID(ctx context.Context, id string) error {
	stmt := sqlf.
		DeleteFrom(SIGN_IN_CODE_MODEL_TABLE).
		Where("id = ?", id)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

func (self *SignInCodeRepositoryImpl) DeleteByExpiresAt(ctx context.Context, expiresAt time.Time) error {
	stmt := sqlf.
		DeleteFrom(SIGN_IN_CODE_MODEL_TABLE).
		Where("expires_at <= ?", expiresAt)

	_, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	return nil
}
