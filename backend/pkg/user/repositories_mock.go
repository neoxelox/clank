package user

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func NewUserRepositoryMock() *UserRepositoryMock {
	return &UserRepositoryMock{}
}

func (m *UserRepositoryMock) Create(ctx context.Context, user User) (*User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *UserRepositoryMock) GetByID(ctx context.Context, id string) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *UserRepositoryMock) GetByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *UserRepositoryMock) GetByOrganizationIDAndEmail(ctx context.Context, organizationID string, email string) (*User, error) {
	args := m.Called(ctx, organizationID, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *UserRepositoryMock) ListByOrganizationID(ctx context.Context, organizationID string) ([]User, error) {
	args := m.Called(ctx, organizationID)
	return args.Get(0).([]User), args.Error(1)
}

func (m *UserRepositoryMock) CountNotDeletedByOrganizationID(ctx context.Context, organizationID string) (int, error) {
	args := m.Called(ctx, organizationID)
	return args.Int(0), args.Error(1)
}

func (m *UserRepositoryMock) UpdateProfile(ctx context.Context, user User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) UpdateEmail(ctx context.Context, id string, email string) error {
	args := m.Called(ctx, id, email)
	return args.Error(0)
}

func (m *UserRepositoryMock) UpdateSettings(ctx context.Context, user User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) UpdateRole(ctx context.Context, id string, role string) error {
	args := m.Called(ctx, id, role)
	return args.Error(0)
}

func (m *UserRepositoryMock) UpdateDeletedAt(ctx context.Context, id string, deletedAt time.Time) error {
	args := m.Called(ctx, id, deletedAt)
	return args.Error(0)
}

type InvitationRepositoryMock struct {
	mock.Mock
}

func NewInvitationRepositoryMock() *InvitationRepositoryMock {
	return &InvitationRepositoryMock{}
}

func (m *InvitationRepositoryMock) Create(ctx context.Context, invitation Invitation) (*Invitation, error) {
	args := m.Called(ctx, invitation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Invitation), args.Error(1)
}

func (m *InvitationRepositoryMock) GetByID(ctx context.Context, id string) (*Invitation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Invitation), args.Error(1)
}

func (m *InvitationRepositoryMock) GetByEmail(ctx context.Context, email string) (*Invitation, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Invitation), args.Error(1)
}

func (m *InvitationRepositoryMock) ListByOrganizationID(ctx context.Context, organizationID string) ([]Invitation, error) {
	args := m.Called(ctx, organizationID)
	return args.Get(0).([]Invitation), args.Error(1)
}

func (m *InvitationRepositoryMock) CountNotExpiredByOrganizationID(ctx context.Context, organizationID string) (int, error) {
	args := m.Called(ctx, organizationID)
	return args.Int(0), args.Error(1)
}

func (m *InvitationRepositoryMock) DeleteByID(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *InvitationRepositoryMock) DeleteByExpiresAt(ctx context.Context, expiresAt time.Time) error {
	args := m.Called(ctx, expiresAt)
	return args.Error(0)
}
