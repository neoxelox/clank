package exporter

import (
	"fmt"
	"time"

	"github.com/neoxelox/kit/util"
)

const (
	ExporterTypeSlack = "SLACK"
	ExporterTypeJira  = "JIRA"
)

func IsExporterType(value string) bool {
	return value == ExporterTypeSlack ||
		value == ExporterTypeJira
}

type ExporterSettings struct {
}

type ExporterJobdata struct {
}

type Exporter struct {
	ID        string
	ProductID string
	Type      string
	Settings  any
	Jobdata   any
	CreatedAt time.Time
	DeletedAt *time.Time
}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (self Exporter) String() string {
	return fmt.Sprintf("<Exporter: %s (%s)>", self.Type, self.ID)
}

func (self Exporter) Equals(other Exporter) bool {
	return util.Equals(self, other)
}

func (self Exporter) Copy() *Exporter {
	return util.Copy(self)
}
