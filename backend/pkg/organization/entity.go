package organization

import (
	"fmt"
	"time"

	"github.com/neoxelox/kit/util"
	emailproviders "gomodules.xyz/email-providers"
)

const (
	ORGANIZATION_DEFAULT_PICTURE           = "https://clank.so/images/pictures/organization.png"
	ORGANIZATION_PLAN_STARTER_MAX_PRODUCTS = 5
	ORGANIZATION_PLAN_TRIAL_MAX_MEMBERS    = 3 + 1
	ORGANIZATION_PLAN_TRIAL_PERIOD         = 14 * 24 * time.Hour
)

func IsSafeDomain(domain string) bool {
	return emailproviders.IsWorkEmail("work@"+domain) && domain != "privaterelay.appleid.com"
}

const (
	OrganizationPlanEnterprise = "ENTERPRISE"
	OrganizationPlanBusiness   = "BUSINESS"
	OrganizationPlanStarter    = "STARTER"
	OrganizationPlanTrial      = "TRIAL"
	OrganizationPlanDemo       = "DEMO"
)

func IsOrganizationPlan(value string) bool {
	return value == OrganizationPlanEnterprise ||
		value == OrganizationPlanBusiness ||
		value == OrganizationPlanStarter ||
		value == OrganizationPlanTrial ||
		value == OrganizationPlanDemo
}

var OrganizationPlanIncludedCapacity = map[string]int{
	OrganizationPlanEnterprise: 1000000, // This is a hard upper limit just in case
	OrganizationPlanBusiness:   100000,
	OrganizationPlanStarter:    10000,
	OrganizationPlanTrial:      100,
	OrganizationPlanDemo:       0,
}

func PlanTrialCapacity(domain string) (int, int) {
	if !IsSafeDomain(domain) {
		return 0, OrganizationPlanIncludedCapacity[OrganizationPlanTrial]
	}

	return OrganizationPlanIncludedCapacity[OrganizationPlanTrial], 0
}

type OrganizationSettings struct {
	DomainSignIn bool
}

func IsDomainSignInSupported(domain string) bool {
	return IsSafeDomain(domain)
}

type OrganizationCapacity struct {
	Included int
	Extra    int
}

type OrganizationUsage struct {
	Value          int
	LastComputedAt time.Time
}

type Organization struct {
	ID          string
	Name        string
	Picture     string
	Domain      string
	Settings    OrganizationSettings
	Plan        string
	TrialEndsAt time.Time
	Capacity    OrganizationCapacity
	Usage       OrganizationUsage
	CreatedAt   time.Time
	DeletedAt   *time.Time
}

func NewOrganization() *Organization {
	return &Organization{}
}

func (self Organization) UsageLeft() int {
	return max(0, (self.Capacity.Included+self.Capacity.Extra)-self.Usage.Value)
}

func (self Organization) String() string {
	return fmt.Sprintf("<Organization: %s (%s)>", self.Name, self.ID)
}

func (self Organization) Equals(other Organization) bool {
	return util.Equals(self, other)
}

func (self Organization) Copy() *Organization {
	return util.Copy(self)
}
