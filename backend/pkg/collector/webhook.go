package collector

import (
	"backend/pkg/config"

	"github.com/neoxelox/kit"
)

const (
	WEBHOOK_COLLECTOR_CALLBACK_SECRET_LENGTH = 32
)

type WebhookCollectorSettings struct {
	CollectorSettings
	APIKey string
}

type WebhookCollectorJobdata struct {
	CollectorJobdata
}

type WebhookCollector struct {
	config   config.Config
	observer *kit.Observer
}

func NewWebhookCollector(observer *kit.Observer, config config.Config) *WebhookCollector {
	return &WebhookCollector{
		config:   config,
		observer: observer,
	}
}
