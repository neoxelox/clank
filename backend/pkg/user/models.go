package user

import (
	"encoding/json"
	"time"
)

const (
	USER_MODEL_TABLE = "\"user\""
)

type UserModel struct {
	ID             string     `db:"id"`
	OrganizationID string     `db:"organization_id"`
	Name           string     `db:"name"`
	Picture        string     `db:"picture"`
	Email          string     `db:"email"`
	Role           string     `db:"role"`
	Settings       []byte     `db:"settings"`
	CreatedAt      time.Time  `db:"created_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
}

func NewUserModel(user User) *UserModel {
	settings, err := json.Marshal(user.Settings)
	if err != nil {
		panic(err)
	}

	return &UserModel{
		ID:             user.ID,
		OrganizationID: user.OrganizationID,
		Name:           user.Name,
		Picture:        user.Picture,
		Email:          user.Email,
		Role:           user.Role,
		Settings:       settings,
		CreatedAt:      user.CreatedAt,
		DeletedAt:      user.DeletedAt,
	}
}

func (self *UserModel) ToEntity() *User {
	var settings UserSettings
	err := json.Unmarshal(self.Settings, &settings)
	if err != nil {
		panic(err)
	}

	return &User{
		ID:             self.ID,
		OrganizationID: self.OrganizationID,
		Name:           self.Name,
		Picture:        self.Picture,
		Email:          self.Email,
		Role:           self.Role,
		Settings:       settings,
		CreatedAt:      self.CreatedAt,
		DeletedAt:      self.DeletedAt,
	}
}

const (
	INVITATION_MODEL_TABLE = "\"invitation\""
)

type InvitationModel struct {
	ID             string    `db:"id"`
	OrganizationID string    `db:"organization_id"`
	Email          string    `db:"email"`
	Role           string    `db:"role"`
	ExpiresAt      time.Time `db:"expires_at"`
}

func NewInvitationModel(invitation Invitation) *InvitationModel {
	return &InvitationModel{
		ID:             invitation.ID,
		OrganizationID: invitation.OrganizationID,
		Email:          invitation.Email,
		Role:           invitation.Role,
		ExpiresAt:      invitation.ExpiresAt,
	}
}

func (self *InvitationModel) ToEntity() *Invitation {
	return &Invitation{
		ID:             self.ID,
		OrganizationID: self.OrganizationID,
		Email:          self.Email,
		Role:           self.Role,
		ExpiresAt:      self.ExpiresAt,
	}
}
