package collector

import (
	"context"
	"time"

	"github.com/leporo/sqlf"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
)

type CollectorRepository struct {
	config   config.Config
	observer *kit.Observer
	database *kit.Database
}

func NewCollectorRepository(observer *kit.Observer, database *kit.Database, config config.Config) *CollectorRepository {
	return &CollectorRepository{
		config:   config,
		observer: observer,
		database: database,
	}
}

func (self *CollectorRepository) Create(ctx context.Context, collector Collector) (*Collector, error) {
	c := NewCollectorModel(collector)

	stmt := sqlf.
		InsertInto(COLLECTOR_MODEL_TABLE).
		Set("id", c.ID).
		Set("product_id", c.ProductID).
		Set("type", c.Type).
		Set("settings", c.Settings).
		Set("jobdata", c.Jobdata).
		Set("created_at", c.CreatedAt).
		Set("deleted_at", c.DeletedAt).
		Returning("*").To(&c)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	return c.ToEntity(), nil
}

func (self *CollectorRepository) GetByID(ctx context.Context, id string) (*Collector, error) {
	var c CollectorModel

	stmt := sqlf.
		Select("*").To(&c).
		From(COLLECTOR_MODEL_TABLE).
		Where("id = ?", id)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return c.ToEntity(), nil
}

func (self *CollectorRepository) ListByProductIDAndType(ctx context.Context,
	productID string, _type string) ([]Collector, error) {
	var cs []CollectorModel

	stmt := sqlf.
		Select("*").To(&cs).
		From(COLLECTOR_MODEL_TABLE).
		Where("product_id = ?", productID).
		Where("type = ?", _type)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return []Collector{}, nil
		}

		return nil, err
	}

	entities := make([]Collector, 0, len(cs))
	for _, c := range cs {
		entities = append(entities, *c.ToEntity())
	}

	return entities, nil
}

func (self *CollectorRepository) ListByProductID(ctx context.Context, productID string) ([]Collector, error) {
	var cs []CollectorModel

	stmt := sqlf.
		Select("*").To(&cs).
		From(COLLECTOR_MODEL_TABLE).
		Where("product_id = ?", productID)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return []Collector{}, nil
		}

		return nil, err
	}

	entities := make([]Collector, 0, len(cs))
	for _, c := range cs {
		entities = append(entities, *c.ToEntity())
	}

	return entities, nil
}

func (self *CollectorRepository) ListIDsByTypeNotDeleted(ctx context.Context, _type string) ([]string, error) {
	var result []struct {
		ID string `db:"id"`
	}

	stmt := sqlf.
		Select("id").To(&result).
		From(COLLECTOR_MODEL_TABLE).
		Where("type = ?", _type).
		Where("deleted_at IS NULL")

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

func (self *CollectorRepository) ExistsByOrganizationID(ctx context.Context, organizationID string) (bool, error) {
	var e bool

	stmt := sqlf.
		Select("").To(&e).
		SubQuery("EXISTS (", ")", sqlf.
			Select("TRUE").
			From(`"collector"`).
			Join(`"product"`, `"collector".product_id = "product".id`).
			Where(`"product".organization_id = ?`, organizationID))

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return false, nil
		}

		return false, err
	}

	return e, nil
}

func (self *CollectorRepository) UpdateSettings(ctx context.Context, collector Collector) error {
	c := NewCollectorModel(collector)

	stmt := sqlf.
		Update(COLLECTOR_MODEL_TABLE).
		Set("settings", c.Settings).
		Where("id = ?", c.ID)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

func (self *CollectorRepository) UpdateSettingsAndDeleted(ctx context.Context, collector Collector) error {
	c := NewCollectorModel(collector)

	stmt := sqlf.
		Update(COLLECTOR_MODEL_TABLE).
		Set("settings", c.Settings).
		Set("deleted_at", c.DeletedAt).
		Where("id = ?", c.ID)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

func (self *CollectorRepository) UpdateJobdata(ctx context.Context, collector Collector) error {
	c := NewCollectorModel(collector)

	stmt := sqlf.
		Update(COLLECTOR_MODEL_TABLE).
		Set("jobdata", c.Jobdata).
		Where("id = ?", c.ID)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

func (self *CollectorRepository) UpdateDeletedAt(ctx context.Context, id string, deletedAt time.Time) error {
	stmt := sqlf.
		Update(COLLECTOR_MODEL_TABLE).
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
