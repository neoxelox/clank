package collector

import "encoding/json"

type TrustpilotCollectorPayloadSettings struct {
	CollectorPayloadSettings
	Domain string `json:"domain"`
}

type PlayStoreCollectorPayloadSettings struct {
	CollectorPayloadSettings
	AppID string `json:"app_id"`
}

type AppStoreCollectorPayloadSettings struct {
	CollectorPayloadSettings
	AppID string `json:"app_id"`
}

type AmazonCollectorPayloadSettings struct {
	CollectorPayloadSettings
	ASIN string `json:"asin"`
}

type IAgoraCollectorPayloadSettings struct {
	CollectorPayloadSettings
	Institution string `json:"institution"`
}

type WebhookCollectorPayloadSettings struct {
	CollectorPayloadSettings
	APIKey string `json:"api_key"`
}

type WidgetCollectorPayloadSettings struct {
	CollectorPayloadSettings
	ClientKey string `json:"client_key"`
}

type CollectorPayloadSettings struct {
}

type CollectorPayload struct {
	ID        string          `json:"id"`
	ProductID string          `json:"product_id"`
	Type      string          `json:"type"`
	Settings  json.RawMessage `json:"settings"`
}

func NewCollectorPayload(collector Collector) *CollectorPayload {
	var err error

	var settings json.RawMessage
	switch collector.Type {
	case CollectorTypeTrustpilot:
		_settings := collector.Settings.(TrustpilotCollectorSettings) // nolint: errcheck
		settings, err = json.Marshal(TrustpilotCollectorPayloadSettings{
			Domain: _settings.Domain,
		})
		if err != nil {
			panic(err)
		}

	case CollectorTypePlayStore:
		_settings := collector.Settings.(PlayStoreCollectorSettings) // nolint: errcheck
		settings, err = json.Marshal(PlayStoreCollectorPayloadSettings{
			AppID: _settings.AppID,
		})
		if err != nil {
			panic(err)
		}

	case CollectorTypeAppStore:
		_settings := collector.Settings.(AppStoreCollectorSettings) // nolint: errcheck
		settings, err = json.Marshal(AppStoreCollectorPayloadSettings{
			AppID: _settings.AppID,
		})
		if err != nil {
			panic(err)
		}

	case CollectorTypeAmazon:
		_settings := collector.Settings.(AmazonCollectorSettings) // nolint: errcheck
		settings, err = json.Marshal(AmazonCollectorPayloadSettings{
			ASIN: _settings.ASIN,
		})
		if err != nil {
			panic(err)
		}

	case CollectorTypeIAgora:
		_settings := collector.Settings.(IAgoraCollectorSettings) // nolint: errcheck
		settings, err = json.Marshal(IAgoraCollectorPayloadSettings{
			Institution: _settings.Institution,
		})
		if err != nil {
			panic(err)
		}

	case CollectorTypeWebhook:
		_settings := collector.Settings.(WebhookCollectorSettings) // nolint: errcheck
		settings, err = json.Marshal(WebhookCollectorPayloadSettings{
			APIKey: _settings.APIKey,
		})
		if err != nil {
			panic(err)
		}

	case CollectorTypeWidget:
		_settings := collector.Settings.(WidgetCollectorSettings) // nolint: errcheck
		settings, err = json.Marshal(WidgetCollectorPayloadSettings{
			ClientKey: _settings.ClientKey,
		})
		if err != nil {
			panic(err)
		}

	default:
		panic(collector.Type)
	}

	return &CollectorPayload{
		ID:        collector.ID,
		ProductID: collector.ProductID,
		Type:      collector.Type,
		Settings:  settings,
	}
}
