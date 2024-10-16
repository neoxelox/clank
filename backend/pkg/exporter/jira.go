package exporter

import (
	"backend/pkg/config"

	"github.com/neoxelox/kit"
)

type JiraExporterSettings struct {
	ExporterSettings
	Board string
}

type JiraExporterJobdata struct {
	ExporterJobdata
}

type JiraExporter struct {
	config   config.Config
	observer *kit.Observer
}

func NewJiraExporter(observer *kit.Observer, config config.Config) *JiraExporter {
	return &JiraExporter{
		config:   config,
		observer: observer,
	}
}
