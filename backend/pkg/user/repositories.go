package user

import (
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user User) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByOrganizationIDAndEmail(ctx context.Context, organizationID string, email string) (*User, error)
	ListByOrganizationID(ctx context.Context, organizationID string) ([]User, error)
	CountNotDeletedByOrganizationID(ctx context.Context, organizationID string) (int, error)
	UpdateProfile(ctx context.Context, user User) error
	UpdateEmail(ctx context.Context, id string, email string) error
	UpdateSettings(ctx context.Context, user User) error
	UpdateRole(ctx context.Context, id string, role string) error
	UpdateDeletedAt(ctx context.Context, id string, deletedAt time.Time) error
}

type InvitationRepository interface {
	Create(ctx context.Context, invitation Invitation) (*Invitation, error)
	GetByID(ctx context.Context, id string) (*Invitation, error)
	GetByEmail(ctx context.Context, email string) (*Invitation, error)
	ListByOrganizationID(ctx context.Context, organizationID string) ([]Invitation, error)
	CountNotExpiredByOrganizationID(ctx context.Context, organizationID string) (int, error)
	DeleteByID(ctx context.Context, id string) error
	DeleteByExpiresAt(ctx context.Context, expiresAt time.Time) error
}
