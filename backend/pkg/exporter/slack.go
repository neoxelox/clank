package exporter

import (
	"backend/pkg/config"

	"github.com/neoxelox/kit"
)

type SlackExporterSettings struct {
	ExporterSettings
	Channel string
}

type SlackExporterJobdata struct {
	ExporterJobdata
}

type SlackExporter struct {
	config   config.Config
	observer *kit.Observer
}

func NewSlackExporter(observer *kit.Observer, config config.Config) *SlackExporter {
	return &SlackExporter{
		config:   config,
		observer: observer,
	}
}
