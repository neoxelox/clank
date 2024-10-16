package product

import (
	"context"
	"time"

	"github.com/leporo/sqlf"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
)

type ProductRepository struct {
	config   config.Config
	observer *kit.Observer
	database *kit.Database
}

func NewProductRepository(observer *kit.Observer, database *kit.Database, config config.Config) *ProductRepository {
	return &ProductRepository{
		config:   config,
		observer: observer,
		database: database,
	}
}

func (self *ProductRepository) Create(ctx context.Context, product Product) (*Product, error) {
	p := NewProductModel(product)

	stmt := sqlf.
		InsertInto(PRODUCT_MODEL_TABLE).
		Set("id", p.ID).
		Set("organization_id", p.OrganizationID).
		Set("name", p.Name).
		Set("picture", p.Picture).
		Set("language", p.Language).
		Set("context", p.Context).
		Set("categories", p.Categories).
		Set("release", p.Release).
		Set("settings", p.Settings).
		Set("usage", p.Usage).
		Set("created_at", p.CreatedAt).
		Set("deleted_at", p.DeletedAt).
		Returning("*").To(&p)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	return p.ToEntity(), nil
}

func (self *ProductRepository) GetByID(ctx context.Context, id string) (*Product, error) {
	var p ProductModel

	stmt := sqlf.
		Select("*").To(&p).
		From(PRODUCT_MODEL_TABLE).
		Where("id = ?", id)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return p.ToEntity(), nil
}

func (self *ProductRepository) ListByOrganizationID(ctx context.Context, organizationID string) ([]Product, error) {
	var ps []ProductModel

	stmt := sqlf.
		Select("*").To(&ps).
		From(PRODUCT_MODEL_TABLE).
		Where("organization_id = ?", organizationID)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return []Product{}, nil
		}

		return nil, err
	}

	entities := make([]Product, 0, len(ps))
	for _, p := range ps {
		entities = append(entities, *p.ToEntity())
	}

	return entities, nil
}

func (self *ProductRepository) ExistsByOrganizationID(ctx context.Context, organizationID string) (bool, error) {
	var e bool

	stmt := sqlf.
		Select("").To(&e).
		SubQuery("EXISTS (", ")", sqlf.
			Select("TRUE").
			From(PRODUCT_MODEL_TABLE).
			Where("organization_id = ?", organizationID))

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return false, nil
		}

		return false, err
	}

	return e, nil
}

func (self *ProductRepository) CountNotDeletedByOrganizationID(
	ctx context.Context, organizationID string) (int, error) {
	var c int

	stmt := sqlf.
		Select("COUNT(*)").To(&c).
		From(PRODUCT_MODEL_TABLE).
		Where("organization_id = ?", organizationID).
		Where("deleted_at IS NULL")

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return 0, nil
		}

		return 0, err
	}

	return c, nil
}

func (self *ProductRepository) UpdateInfo(ctx context.Context, product Product) error {
	p := NewProductModel(product)

	stmt := sqlf.
		Update(PRODUCT_MODEL_TABLE).
		Set("name", p.Name).
		Set("picture", p.Picture).
		Set("context", p.Context).
		Set("categories", p.Categories).
		Set("release", p.Release).
		Where("id = ?", p.ID)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

func (self *ProductRepository) UpdateSettings(ctx context.Context, product Product) error {
	p := NewProductModel(product)

	stmt := sqlf.
		Update(PRODUCT_MODEL_TABLE).
		Set("settings", p.Settings).
		Where("id = ?", p.ID)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

func (self *ProductRepository) UpdateDeletedAt(ctx context.Context, id string, deletedAt time.Time) error {
	stmt := sqlf.
		Update(PRODUCT_MODEL_TABLE).
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
func (self *ProductRepository) Delete(ctx context.Context, id string) error {
	now := time.Now()

	err := self.database.Transaction(ctx, nil, func(ctx context.Context) error {
		stmt := sqlf.
			Update(`"collector"`).
			Set("deleted_at", now).
			Where("product_id = ?", id)

		_, err := self.database.Exec(ctx, stmt)
		if err != nil {
			return err
		}

		stmt = sqlf.
			Update(`"exporter"`).
			Set("deleted_at", now).
			Where("product_id = ?", id)

		_, err = self.database.Exec(ctx, stmt)
		if err != nil {
			return err
		}

		stmt = sqlf.
			Update(`"product"`).
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
