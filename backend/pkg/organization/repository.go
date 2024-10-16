package organization

import (
	"context"
	"time"
)

type OrganizationRepository interface {
	Create(ctx context.Context, organization Organization) (*Organization, error)
	GetByID(ctx context.Context, id string) (*Organization, error)
	GetByIDForUpdate(ctx context.Context, id string) (*Organization, error)
	GetByDomain(ctx context.Context, domain string) (*Organization, error)
	ListIDsByNotDeleted(ctx context.Context) ([]string, error)
	CountCollectedFeedbacksThisMonthPerProduct(ctx context.Context, id string) (map[string]int, error)
	UpdateProfile(ctx context.Context, organization Organization) error
	UpdateSettings(ctx context.Context, organization Organization) error
	UpdateDomain(ctx context.Context, id string, domain string) error
	DowngradeTrialsByEndsAt(ctx context.Context, endsAt time.Time) error
	UpdateOrganizationUsage(ctx context.Context, organization Organization) error
	UpdateProductsUsage(ctx context.Context, usagePerProduct map[string]int) error
	UpdateDeletedAt(ctx context.Context, id string, deletedAt time.Time) error
	Delete(ctx context.Context, id string) error
}
