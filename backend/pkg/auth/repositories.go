package auth

import (
	"context"
	"time"
)

type SessionRepository interface {
	Create(ctx context.Context, session Session) (*Session, error)
	GetByID(ctx context.Context, id string) (*Session, error)
	UpdateProvider(ctx context.Context, id string, provider string) error
	UpdateSeen(ctx context.Context, session Session) error
	UpdateExpiredAt(ctx context.Context, id string, expiredAt time.Time) error
}

type SignInCodeRepository interface {
	Create(ctx context.Context, signInCode SignInCode) (*SignInCode, error)
	GetByID(ctx context.Context, id string) (*SignInCode, error)
	GetByEmail(ctx context.Context, email string) (*SignInCode, error)
	UpdateAttempts(ctx context.Context, id string, attempts int) error
	DeleteByID(ctx context.Context, id string) error
	DeleteByExpiresAt(ctx context.Context, expiresAt time.Time) error
}
