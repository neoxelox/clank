package user

import (
	"context"
	"time"

	"github.com/leporo/sqlf"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
)

type UserRepositoryImpl struct {
	config   config.Config
	observer *kit.Observer
	database *kit.Database
}

func NewUserRepositoryImpl(observer *kit.Observer, database *kit.Database, config config.Config) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		config:   config,
		observer: observer,
		database: database,
	}
}

func (self *UserRepositoryImpl) Create(ctx context.Context, user User) (*User, error) {
	u := NewUserModel(user)

	stmt := sqlf.
		InsertInto(USER_MODEL_TABLE).
		Set("id", u.ID).
		Set("organization_id", u.OrganizationID).
		Set("name", u.Name).
		Set("picture", u.Picture).
		Set("email", u.Email).
		Set("role", u.Role).
		Set("settings", u.Settings).
		Set("created_at", u.CreatedAt).
		Set("deleted_at", u.DeletedAt).
		Returning("*").To(&u)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	return u.ToEntity(), nil
}

func (self *UserRepositoryImpl) GetByID(ctx context.Context, id string) (*User, error) {
	var u UserModel

	stmt := sqlf.
		Select("*").To(&u).
		From(USER_MODEL_TABLE).
		Where("id = ?", id)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return u.ToEntity(), nil
}

func (self *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*User, error) {
	var u UserModel

	stmt := sqlf.
		Select("*").To(&u).
		From(USER_MODEL_TABLE).
		Where("email = ?", email)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return u.ToEntity(), nil
}

func (self *UserRepositoryImpl) GetByOrganizationIDAndEmail(ctx context.Context,
	organizationID string, email string) (*User, error) {
	var u UserModel

	stmt := sqlf.
		Select("*").To(&u).
		From(USER_MODEL_TABLE).
		Where("organization_id = ?", organizationID).
		Where("email = ?", email)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return u.ToEntity(), nil
}

func (self *UserRepositoryImpl) ListByOrganizationID(ctx context.Context, organizationID string) ([]User, error) {
	var us []UserModel

	stmt := sqlf.
		Select("*").To(&us).
		From(USER_MODEL_TABLE).
		Where("organization_id = ?", organizationID)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return []User{}, nil
		}

		return nil, err
	}

	entities := make([]User, 0, len(us))
	for _, u := range us {
		entities = append(entities, *u.ToEntity())
	}

	return entities, nil
}

func (self *UserRepositoryImpl) CountNotDeletedByOrganizationID(ctx context.Context, organizationID string) (int, error) {
	var c int

	stmt := sqlf.
		Select("COUNT(*)").To(&c).
		From(USER_MODEL_TABLE).
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

func (self *UserRepositoryImpl) UpdateProfile(ctx context.Context, user User) error {
	u := NewUserModel(user)

	stmt := sqlf.
		Update(USER_MODEL_TABLE).
		Set("name", u.Name).
		Set("picture", u.Picture).
		Where("id = ?", u.ID)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

func (self *UserRepositoryImpl) UpdateEmail(ctx context.Context, id string, email string) error {
	stmt := sqlf.
		Update(USER_MODEL_TABLE).
		Set("email", email).
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

func (self *UserRepositoryImpl) UpdateSettings(ctx context.Context, user User) error {
	u := NewUserModel(user)

	stmt := sqlf.
		Update(USER_MODEL_TABLE).
		Set("settings", u.Settings).
		Where("id = ?", u.ID)

	affected, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	if affected != 1 {
		return kit.ErrDatabaseUnexpectedEffect.Raise(affected, 1)
	}

	return nil
}

func (self *UserRepositoryImpl) UpdateRole(ctx context.Context, id string, role string) error {
	stmt := sqlf.
		Update(USER_MODEL_TABLE).
		Set("role", role).
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

func (self *UserRepositoryImpl) UpdateDeletedAt(ctx context.Context, id string, deletedAt time.Time) error {
	stmt := sqlf.
		Update(USER_MODEL_TABLE).
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

type InvitationRepositoryImpl struct {
	config   config.Config
	observer *kit.Observer
	database *kit.Database
}

func NewInvitationRepositoryImpl(observer *kit.Observer, database *kit.Database,
	config config.Config) *InvitationRepositoryImpl {
	return &InvitationRepositoryImpl{
		config:   config,
		observer: observer,
		database: database,
	}
}

func (self *InvitationRepositoryImpl) Create(ctx context.Context, invitation Invitation) (*Invitation, error) {
	i := NewInvitationModel(invitation)

	stmt := sqlf.
		InsertInto(INVITATION_MODEL_TABLE).
		Set("id", i.ID).
		Set("organization_id", i.OrganizationID).
		Set("email", i.Email).
		Set("role", i.Role).
		Set("expires_at", i.ExpiresAt).
		Returning("*").To(&i)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	return i.ToEntity(), nil
}

func (self *InvitationRepositoryImpl) GetByID(ctx context.Context, id string) (*Invitation, error) {
	var i InvitationModel

	stmt := sqlf.
		Select("*").To(&i).
		From(INVITATION_MODEL_TABLE).
		Where("id = ?", id)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return i.ToEntity(), nil
}

func (self *InvitationRepositoryImpl) GetByEmail(ctx context.Context, email string) (*Invitation, error) {
	var i InvitationModel

	stmt := sqlf.
		Select("*").To(&i).
		From(INVITATION_MODEL_TABLE).
		Where("email = ?", email)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return nil, nil
		}

		return nil, err
	}

	return i.ToEntity(), nil
}

func (self *InvitationRepositoryImpl) ListByOrganizationID(ctx context.Context,
	organizationID string) ([]Invitation, error) {
	var is []InvitationModel

	stmt := sqlf.
		Select("*").To(&is).
		From(INVITATION_MODEL_TABLE).
		Where("organization_id = ?", organizationID)

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return []Invitation{}, nil
		}

		return nil, err
	}

	entities := make([]Invitation, 0, len(is))
	for _, i := range is {
		entities = append(entities, *i.ToEntity())
	}

	return entities, nil
}

func (self *InvitationRepositoryImpl) CountNotExpiredByOrganizationID(ctx context.Context,
	organizationID string) (int, error) {
	var c int

	stmt := sqlf.
		Select("COUNT(*)").To(&c).
		From(INVITATION_MODEL_TABLE).
		Where("organization_id = ?", organizationID).
		Where("expires_at > ?", time.Now())

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return 0, nil
		}

		return 0, err
	}

	return c, nil
}

func (self *InvitationRepositoryImpl) DeleteByID(ctx context.Context, id string) error {
	stmt := sqlf.
		DeleteFrom(INVITATION_MODEL_TABLE).
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

func (self *InvitationRepositoryImpl) DeleteByExpiresAt(ctx context.Context, expiresAt time.Time) error {
	stmt := sqlf.
		DeleteFrom(INVITATION_MODEL_TABLE).
		Where("expires_at <= ?", expiresAt)

	_, err := self.database.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	return nil
}
