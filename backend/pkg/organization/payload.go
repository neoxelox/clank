package organization

import (
	"time"

	"github.com/neoxelox/kit/util"
)

type OrganizationPayloadSettings struct {
	DomainSignIn            bool `json:"domain_sign_in"`
	IsDomainSignInSupported bool `json:"is_domain_sign_in_supported"`
}

type OrganizationPayloadCapacity struct {
	Included int `json:"included"`
	Extra    int `json:"extra"`
}

type OrganizationPayload struct {
	ID          string                      `json:"id"`
	Name        string                      `json:"name"`
	Picture     string                      `json:"picture"`
	Domain      string                      `json:"domain"`
	Settings    OrganizationPayloadSettings `json:"settings"`
	Plan        string                      `json:"plan"`
	TrialEndsAt *time.Time                  `json:"trial_ends_at,omitempty"`
	Capacity    OrganizationPayloadCapacity `json:"capacity"`
	Usage       int                         `json:"usage"`
}

func NewOrganizationPayload(organization Organization) *OrganizationPayload {
	var trialEndsAt *time.Time
	if organization.Plan == OrganizationPlanTrial {
		trialEndsAt = util.Pointer(organization.TrialEndsAt)
	}

	return &OrganizationPayload{
		ID:      organization.ID,
		Name:    organization.Name,
		Picture: organization.Picture,
		Domain:  organization.Domain,
		Settings: OrganizationPayloadSettings{
			DomainSignIn:            organization.Settings.DomainSignIn,
			IsDomainSignInSupported: IsDomainSignInSupported(organization.Domain),
		},
		Plan:        organization.Plan,
		TrialEndsAt: trialEndsAt,
		Capacity: OrganizationPayloadCapacity{
			Included: organization.Capacity.Included,
			Extra:    organization.Capacity.Extra,
		},
		Usage: organization.Usage.Value,
	}
}
