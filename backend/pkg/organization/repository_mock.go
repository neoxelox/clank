package organization

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type OrganizationRepositoryMock struct {
	mock.Mock
}

func NewOrganizationRepositoryMock() *OrganizationRepositoryMock {
	return &OrganizationRepositoryMock{}
}

func (m *OrganizationRepositoryMock) Create(ctx context.Context, organization Organization) (*Organization, error) {
	args := m.Called(ctx, organization)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Organization), args.Error(1)
}

func (m *OrganizationRepositoryMock) GetByID(ctx context.Context, id string) (*Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Organization), args.Error(1)
}

func (m *OrganizationRepositoryMock) GetByIDForUpdate(ctx context.Context, id string) (*Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Organization), args.Error(1)
}

func (m *OrganizationRepositoryMock) GetByDomain(ctx context.Context, domain string) (*Organization, error) {
	args := m.Called(ctx, domain)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Organization), args.Error(1)
}

func (m *OrganizationRepositoryMock) ListIDsByNotDeleted(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *OrganizationRepositoryMock) CountCollectedFeedbacksThisMonthPerProduct(ctx context.Context, id string) (map[string]int, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *OrganizationRepositoryMock) UpdateProfile(ctx context.Context, organization Organization) error {
	args := m.Called(ctx, organization)
	return args.Error(0)
}

func (m *OrganizationRepositoryMock) UpdateSettings(ctx context.Context, organization Organization) error {
	args := m.Called(ctx, organization)
	return args.Error(0)
}

func (m *OrganizationRepositoryMock) UpdateDomain(ctx context.Context, id string, domain string) error {
	args := m.Called(ctx, id, domain)
	return args.Error(0)
}

func (m *OrganizationRepositoryMock) DowngradeTrialsByEndsAt(ctx context.Context, endsAt time.Time) error {
	args := m.Called(ctx, endsAt)
	return args.Error(0)
}

func (m *OrganizationRepositoryMock) UpdateOrganizationUsage(ctx context.Context, organization Organization) error {
	args := m.Called(ctx, organization)
	return args.Error(0)
}

func (m *OrganizationRepositoryMock) UpdateProductsUsage(ctx context.Context, usagePerProduct map[string]int) error {
	args := m.Called(ctx, usagePerProduct)
	return args.Error(0)
}

func (m *OrganizationRepositoryMock) UpdateDeletedAt(ctx context.Context, id string, deletedAt time.Time) error {
	args := m.Called(ctx, id, deletedAt)
	return args.Error(0)
}

func (m *OrganizationRepositoryMock) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
