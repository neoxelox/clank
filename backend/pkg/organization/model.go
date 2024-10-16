package organization

import (
	"encoding/json"
	"time"
)

const (
	ORGANIZATION_MODEL_TABLE = "\"organization\""
)

type OrganizationModel struct {
	ID          string     `db:"id"`
	Name        string     `db:"name"`
	Picture     string     `db:"picture"`
	Domain      string     `db:"domain"`
	Settings    []byte     `db:"settings"`
	Plan        string     `db:"plan"`
	TrialEndsAt time.Time  `db:"trial_ends_at"`
	Capacity    []byte     `db:"capacity"`
	Usage       []byte     `db:"usage"`
	CreatedAt   time.Time  `db:"created_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

func NewOrganizationModel(organization Organization) *OrganizationModel {
	settings, err := json.Marshal(organization.Settings)
	if err != nil {
		panic(err)
	}

	capacity, err := json.Marshal(organization.Capacity)
	if err != nil {
		panic(err)
	}

	usage, err := json.Marshal(organization.Usage)
	if err != nil {
		panic(err)
	}

	return &OrganizationModel{
		ID:          organization.ID,
		Name:        organization.Name,
		Picture:     organization.Picture,
		Domain:      organization.Domain,
		Settings:    settings,
		Plan:        organization.Plan,
		TrialEndsAt: organization.TrialEndsAt,
		Capacity:    capacity,
		Usage:       usage,
		CreatedAt:   organization.CreatedAt,
		DeletedAt:   organization.DeletedAt,
	}
}

func (self *OrganizationModel) ToEntity() *Organization {
	var settings OrganizationSettings
	err := json.Unmarshal(self.Settings, &settings)
	if err != nil {
		panic(err)
	}

	var capacity OrganizationCapacity
	err = json.Unmarshal(self.Capacity, &capacity)
	if err != nil {
		panic(err)
	}

	var usage OrganizationUsage
	err = json.Unmarshal(self.Usage, &usage)
	if err != nil {
		panic(err)
	}

	return &Organization{
		ID:          self.ID,
		Name:        self.Name,
		Picture:     self.Picture,
		Domain:      self.Domain,
		Settings:    settings,
		Plan:        self.Plan,
		TrialEndsAt: self.TrialEndsAt,
		Capacity:    capacity,
		Usage:       usage,
		CreatedAt:   self.CreatedAt,
		DeletedAt:   self.DeletedAt,
	}
}
