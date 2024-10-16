package collector

import (
	"backend/pkg/config"

	"github.com/neoxelox/kit"
)

const (
	WIDGET_COLLECTOR_CALLBACK_SECRET_LENGTH = 32
)

type WidgetCollectorSettings struct {
	CollectorSettings
	ClientKey string
}

type WidgetCollectorJobdata struct {
	CollectorJobdata
}

type WidgetCollector struct {
	config   config.Config
	observer *kit.Observer
}

func NewWidgetCollector(observer *kit.Observer, config config.Config) *WidgetCollector {
	return &WidgetCollector{
		config:   config,
		observer: observer,
	}
}
