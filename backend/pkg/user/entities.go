package user

import (
	"fmt"
	"time"

	"github.com/neoxelox/kit/util"
)

const (
	USER_DEFAULT_PICTURE = "https://clank.so/images/pictures/user.png"
)

const (
	UserRoleAdmin  = "ADMIN"
	UserRoleMember = "MEMBER"
)

func IsUserRole(value string) bool {
	return value == UserRoleAdmin ||
		value == UserRoleMember
}

type UserSettings struct {
}

type User struct {
	ID             string
	OrganizationID string
	Name           string
	Picture        string
	Email          string
	Role           string
	Settings       UserSettings
	CreatedAt      time.Time
	DeletedAt      *time.Time
}

func NewUser() *User {
	return &User{}
}

func (self User) String() string {
	return fmt.Sprintf("<User: %s (%s)>", self.Name, self.ID)
}

func (self User) Equals(other User) bool {
	return util.Equals(self, other)
}

func (self User) Copy() *User {
	return util.Copy(self)
}

const (
	INVITATION_EXPIRATION = 7 * 24 * time.Hour
)

type Invitation struct {
	ID             string
	OrganizationID string
	Email          string
	Role           string
	ExpiresAt      time.Time
}

func NewInvitation() *Invitation {
	return &Invitation{}
}

func (self Invitation) String() string {
	return fmt.Sprintf("<Invitation: %s (%s)>", self.Email, self.ID)
}

func (self Invitation) Equals(other Invitation) bool {
	return util.Equals(self, other)
}

func (self Invitation) Copy() *Invitation {
	return util.Copy(self)
}
