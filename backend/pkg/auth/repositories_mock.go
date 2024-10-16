package auth

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type SessionRepositoryMock struct {
	mock.Mock
}

func NewSessionRepositoryMock() *SessionRepositoryMock {
	return &SessionRepositoryMock{}
}

func (m *SessionRepositoryMock) Create(ctx context.Context, session Session) (*Session, error) {
	args := m.Called(ctx, session)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Session), args.Error(1)
}

func (m *SessionRepositoryMock) GetByID(ctx context.Context, id string) (*Session, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Session), args.Error(1)
}

func (m *SessionRepositoryMock) UpdateProvider(ctx context.Context, id string, provider string) error {
	args := m.Called(ctx, id, provider)
	return args.Error(0)
}

func (m *SessionRepositoryMock) UpdateSeen(ctx context.Context, session Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *SessionRepositoryMock) UpdateExpiredAt(ctx context.Context, id string, expiredAt time.Time) error {
	args := m.Called(ctx, id, expiredAt)
	return args.Error(0)
}

type SignInCodeRepositoryMock struct {
	mock.Mock
}

func NewSignInCodeRepositoryMock() *SignInCodeRepositoryMock {
	return &SignInCodeRepositoryMock{}
}

func (m *SignInCodeRepositoryMock) Create(ctx context.Context, signInCode SignInCode) (*SignInCode, error) {
	args := m.Called(ctx, signInCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SignInCode), args.Error(1)
}

func (m *SignInCodeRepositoryMock) GetByID(ctx context.Context, id string) (*SignInCode, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SignInCode), args.Error(1)
}

func (m *SignInCodeRepositoryMock) GetByEmail(ctx context.Context, email string) (*SignInCode, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SignInCode), args.Error(1)
}

func (m *SignInCodeRepositoryMock) UpdateAttempts(ctx context.Context, id string, attempts int) error {
	args := m.Called(ctx, id, attempts)
	return args.Error(0)
}

func (m *SignInCodeRepositoryMock) DeleteByID(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *SignInCodeRepositoryMock) DeleteByExpiresAt(ctx context.Context, expiresAt time.Time) error {
	args := m.Called(ctx, expiresAt)
	return args.Error(0)
}
