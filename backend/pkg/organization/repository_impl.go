package organization

import (
	"context"
	"time"

	"github.com/leporo/sqlf"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
	"backend/pkg/util"
)

type OrganizationRepositoryImpl struct {
	config   config.Config
	observer *kit.Observer
	database *kit.Database
}

func NewOrganizationRepositoryImpl(observer *kit.Observer, database *kit.Database,
	config config.Config) *OrganizationRepositoryImpl {
	return &OrganizationRepositoryImpl{
		config:   config,
		observer: observer,
		database: database,
	}
}

func (self *OrganizationRepositoryImpl) Create(ctx context.Context, organization Organization) (*Organization, error) {
	o := NewOrganizationModel(organization)

	stmt := sqlf.
		InsertInto(ORGANIZATION_MODEL_TABLE).
		Set("id", o.ID).
		Set("name", o.Name).
		Set("picture", o.Picture).
		Set("domain", o.Domain).
		Set("settings", o.Settings).
		Set("plan", o.Plan).
		Set("trial_ends_at", o.TrialEndsAt).
		Set("capacity", o.Capacity).
		Set("usage", o.Usage).
		Set("created_at", o.CreatedAt).
		Set("deleted_at", o.DeletedAt).
		Returning("*").To(&o)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	return o.ToEntity(), nil
}

func (self *OrganizationRepositoryImpl) GetByID(ctx context.Context, id string) (*Organization, error) {
	var o OrganizationModel

	stmt := sqlf.
		Select("*").To(&o).
		From(ORGANIZATION_MODEL_TABLE).
		Where("id = ?", id)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return o.ToEntity(), nil
}

func (self *OrganizationRepositoryImpl) GetByIDForUpdate(ctx context.Context, id string) (*Organization, error) {
	var o OrganizationModel

	stmt := sqlf.
		Select("*").To(&o).
		From(ORGANIZATION_MODEL_TABLE).
		Where("id = ?", id).
		Clause("FOR NO KEY UPDATE NOWAIT")

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return o.ToEntity(), nil
}

func (self *OrganizationRepositoryImpl) GetByDomain(ctx context.Context, domain string) (*Organization, error) {
	var o OrganizationModel

	stmt := sqlf.
		Select("*").To(&o).
		From(ORGANIZATION_MODEL_TABLE).
		Where("domain = ?", domain)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return o.ToEntity(), nil
}

func (self *OrganizationRepositoryImpl) ListIDsByNotDeleted(ctx context.Context) ([]string, error) {
	var result []struct {
		ID string `db:"id"`
	}

	stmt := sqlf.
		Select("id").To(&result).
		From(ORGANIZATION_MODEL_TABLE).
		Where("deleted_at is NULL")

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return []string{}, nil
		}

		return nil, err
	}

	ids := make([]string, 0, len(result))
	for _, res := range result {
		ids = append(ids, res.ID)
	}

	return ids, nil
}

// This is not very beautiful
func (self *OrganizationRepositoryImpl) CountCollectedFeedbacksThisMonthPerProduct(ctx context.Context, id string) (map[string]int, error) {
	var result []struct {
		ID    string `db:"id"`
		Count int    `db:"count"`
	}
	count := make(map[string]int)

	firstDayOfMonth := util.StartOfDay(util.StartOfMonth(time.Now()))

	stmt := sqlf.
		Select(`"product".id, COUNT(CASE WHEN "feedback".collected_at >= ? THEN 1 ELSE NULL END)`, firstDayOfMonth).To(&result).
		From(`"product"`).
		LeftJoin(`"feedback"`, `"feedback".product_id = "product".id`).
		Where(`"product".organization_id = ?`, id).
		Where(`("product".deleted_at IS NULL OR "product".deleted_at >= ?)`, firstDayOfMonth).
		GroupBy(`"product".id`)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return count, nil
		}

		return nil, err
	}

	for _, res := range result {
		count[res.ID] = res.Count
	}

	return count, nil
}

func (self *OrganizationRepositoryImpl) UpdateProfile(ctx context.Context, organization Organization) error {
	o := NewOrganizationModel(organization)

	stmt := sqlf.
		Update(ORGANIZATION_MODEL_TABLE).
		Set("name", o.Name).
		Set("picture", o.Picture).
		Where("id = ?", o.ID)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

func (self *OrganizationRepositoryImpl) UpdateSettings(ctx context.Context, organization Organization) error {
	o := NewOrganizationModel(organization)

	stmt := sqlf.
		Update(ORGANIZATION_MODEL_TABLE).
		Set("settings", o.Settings).
		Where("id = ?", o.ID)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

func (self *OrganizationRepositoryImpl) UpdateDomain(ctx context.Context, id string, domain string) error {
	stmt := sqlf.
		Update(ORGANIZATION_MODEL_TABLE).
		Set("domain", domain).
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

func (self *OrganizationRepositoryImpl) DowngradeTrialsByEndsAt(ctx context.Context, endsAt time.Time) error {
	stmt := sqlf.
		Update(ORGANIZATION_MODEL_TABLE).
		Set("plan", OrganizationPlanDemo).
		Set("capacity['Included']", OrganizationPlanIncludedCapacity[OrganizationPlanDemo]).
		Where("plan = ?", OrganizationPlanTrial).
		Where("trial_ends_at <= ?", endsAt)

	_, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	return nil
}

func (self *OrganizationRepositoryImpl) UpdateOrganizationUsage(ctx context.Context, organization Organization) error {
	o := NewOrganizationModel(organization)

	stmt := sqlf.
		Update(ORGANIZATION_MODEL_TABLE).
		Set("capacity", o.Capacity).
		Set("usage", o.Usage).
		Where("id = ?", o.ID)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

// This is not very beautiful
func (self *OrganizationRepositoryImpl) UpdateProductsUsage(ctx context.Context, usagePerProduct map[string]int) error {
	products := len(usagePerProduct)
	productIDs := make([]string, 0, products)
	productUsages := make([]int, 0, products)
	for id, usage := range usagePerProduct {
		productIDs = append(productIDs, id)
		productUsages = append(productUsages, usage)
	}

	stmt := sqlf.
		Update(`"product"`).
		Clause("SET usage = bulk_usage").
		From("").
		SubQuery("(", ")", sqlf.
			Select("UNNEST(?::varchar[]) as bulk_id", productIDs).
			Select("UNNEST(?::bigint[]) as bulk_usage", productUsages)).
		Where("id = bulk_id")

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != products {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, products)
	}

	return nil
}

func (self *OrganizationRepositoryImpl) UpdateDeletedAt(ctx context.Context, id string, deletedAt time.Time) error {
	stmt := sqlf.
		Update(ORGANIZATION_MODEL_TABLE).
		Set("deleted_at", deletedAt).
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

// This is not very beautiful
func (self *OrganizationRepositoryImpl) Delete(ctx context.Context, id string) error {
	now := time.Now()

	err := self.database.Transaction(ctx, nil, func(ctx context.Context) error {
		var productIDsRaw []struct {
			ID string `db:"id"`
		}

		stmt := sqlf.
			Select("id").To(&productIDsRaw).
			From(`"product"`).
			Where("organization_id = ?", id)

		err := self.database.Query(ctx, stmt)
		if err != nil && !kit.ErrDatabaseNoRows.Is(err) {
			return err
		}

		productIDs := make([]any, 0, len(productIDsRaw))
		for _, raw := range productIDsRaw {
			productIDs = append(productIDs, raw.ID)
		}

		if len(productIDs) > 0 {
			stmt = sqlf.
				Update(`"collector"`).
				Set("deleted_at", now).
				Where("product_id").In(productIDs...)

			_, err = self.database.Exec(ctx, stmt)
			if err != nil {
				return err
			}

			stmt = sqlf.
				Update(`"exporter"`).
				Set("deleted_at", now).
				Where("product_id").In(productIDs...)

			_, err = self.database.Exec(ctx, stmt)
			if err != nil {
				return err
			}

			stmt = sqlf.
				Update(`"product"`).
				Set("deleted_at", now).
				Where("id").In(productIDs...)

			_, err = self.database.Exec(ctx, stmt)
			if err != nil {
				return err
			}
		}

		stmt = sqlf.
			DeleteFrom(`"invitation"`).
			Where("organization_id = ?", id)

		_, err = self.database.Exec(ctx, stmt)
		if err != nil {
			return err
		}

		stmt = sqlf.
			Update(`"user"`).
			Set("deleted_at", now).
			Where("organization_id = ?", id)

		_, err = self.database.Exec(ctx, stmt)
		if err != nil {
			return err
		}

		stmt = sqlf.
			Update(`"organization"`).
			Set("deleted_at", now).
			Where("id = ?", id)

		_, err = self.database.Exec(ctx, stmt)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
