package exporter

import (
	"context"
	"time"

	"github.com/leporo/sqlf"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
)

type ExporterRepository struct {
	config   config.Config
	observer *kit.Observer
	database *kit.Database
}

func NewExporterRepository(observer *kit.Observer, database *kit.Database, config config.Config) *ExporterRepository {
	return &ExporterRepository{
		config:   config,
		observer: observer,
		database: database,
	}
}

func (self *ExporterRepository) Create(ctx context.Context, exporter Exporter) (*Exporter, error) {
	e := NewExporterModel(exporter)

	stmt := sqlf.
		InsertInto(EXPORTER_MODEL_TABLE).
		Set("id", e.ID).
		Set("product_id", e.ProductID).
		Set("type", e.Type).
		Set("settings", e.Settings).
		Set("jobdata", e.Jobdata).
		Set("created_at", e.CreatedAt).
		Set("deleted_at", e.DeletedAt).
		Returning("*").To(&e)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	return e.ToEntity(), nil
}

func (self *ExporterRepository) GetByID(ctx context.Context, id string) (*Exporter, error) {
	var e ExporterModel

	stmt := sqlf.
		Select("*").To(&e).
		From(EXPORTER_MODEL_TABLE).
		Where("id = ?", id)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return e.ToEntity(), nil
}

func (self *ExporterRepository) ListByProductIDAndType(ctx context.Context,
	productID string, _type string) ([]Exporter, error) {
	var es []ExporterModel

	stmt := sqlf.
		Select("*").To(&es).
		From(EXPORTER_MODEL_TABLE).
		Where("product_id = ?", productID).
		Where("type = ?", _type)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return []Exporter{}, nil
		}

		return nil, err
	}

	entities := make([]Exporter, 0, len(es))
	for _, e := range es {
		entities = append(entities, *e.ToEntity())
	}

	return entities, nil
}

func (self *ExporterRepository) ListByProductID(ctx context.Context, productID string) ([]Exporter, error) {
	var es []ExporterModel

	stmt := sqlf.
		Select("*").To(&es).
		From(EXPORTER_MODEL_TABLE).
		Where("product_id = ?", productID)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return []Exporter{}, nil
		}

		return nil, err
	}

	entities := make([]Exporter, 0, len(es))
	for _, e := range es {
		entities = append(entities, *e.ToEntity())
	}

	return entities, nil
}

func (self *ExporterRepository) ExistsByOrganizationID(ctx context.Context, organizationID string) (bool, error) {
	var e bool

	stmt := sqlf.
		Select("").To(&e).
		SubQuery("EXISTS (", ")", sqlf.
			Select("TRUE").
			From(`"exporter"`).
			Join(`"product"`, `"exporter".product_id = "product".id`).
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

func (self *ExporterRepository) UpdateSettings(ctx context.Context, exporter Exporter) error {
	e := NewExporterModel(exporter)

	stmt := sqlf.
		Update(EXPORTER_MODEL_TABLE).
		Set("settings", e.Settings).
		Where("id = ?", e.ID)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

func (self *ExporterRepository) UpdateSettingsAndDeleted(ctx context.Context, exporter Exporter) error {
	e := NewExporterModel(exporter)

	stmt := sqlf.
		Update(EXPORTER_MODEL_TABLE).
		Set("settings", e.Settings).
		Set("deleted_at", e.DeletedAt).
		Where("id = ?", e.ID)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

func (self *ExporterRepository) UpdateDeletedAt(ctx context.Context, id string, deletedAt time.Time) error {
	stmt := sqlf.
		Update(EXPORTER_MODEL_TABLE).
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
