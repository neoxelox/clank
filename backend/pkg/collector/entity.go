package collector

import (
	"fmt"
	"time"

	"github.com/neoxelox/kit/util"
)

const (
	CollectorTypeTrustpilot = "TRUSTPILOT"
	CollectorTypePlayStore  = "PLAY_STORE"
	CollectorTypeAppStore   = "APP_STORE"
	CollectorTypeAmazon     = "AMAZON"
	CollectorTypeIAgora     = "IAGORA"
	CollectorTypeWebhook    = "WEBHOOK"
	CollectorTypeWidget     = "WIDGET"
)

func IsCollectorType(value string) bool {
	return value == CollectorTypeTrustpilot ||
		value == CollectorTypePlayStore ||
		value == CollectorTypeAppStore ||
		value == CollectorTypeAmazon ||
		value == CollectorTypeIAgora ||
		value == CollectorTypeWebhook ||
		value == CollectorTypeWidget
}

type CollectorSettings struct {
}

type CollectorJobdata struct {
}

type Collector struct {
	ID        string
	ProductID string
	Type      string
	Settings  any
	Jobdata   any
	CreatedAt time.Time
	DeletedAt *time.Time
}

func NewCollector() *Collector {
	return &Collector{}
}

func (self Collector) String() string {
	return fmt.Sprintf("<Collector: %s (%s)>", self.Type, self.ID)
}

func (self Collector) Equals(other Collector) bool {
	return util.Equals(self, other)
}

func (self Collector) Copy() *Collector {
	return util.Copy(self)
}
