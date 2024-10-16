package collector

import (
	"encoding/json"
	"time"
)

const (
	COLLECTOR_MODEL_TABLE = "\"collector\""
)

type CollectorModel struct {
	ID        string     `db:"id"`
	ProductID string     `db:"product_id"`
	Type      string     `db:"type"`
	Settings  []byte     `db:"settings"`
	Jobdata   []byte     `db:"jobdata"`
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func NewCollectorModel(collector Collector) *CollectorModel {
	settings, err := json.Marshal(collector.Settings)
	if err != nil {
		panic(err)
	}

	jobdata, err := json.Marshal(collector.Jobdata)
	if err != nil {
		panic(err)
	}

	return &CollectorModel{
		ID:        collector.ID,
		ProductID: collector.ProductID,
		Type:      collector.Type,
		Settings:  settings,
		Jobdata:   jobdata,
		CreatedAt: collector.CreatedAt,
		DeletedAt: collector.DeletedAt,
	}
}

func (self *CollectorModel) ToEntity() *Collector {
	var settings any
	var jobdata any
	switch self.Type {
	case CollectorTypeTrustpilot:
		var _settings TrustpilotCollectorSettings
		err := json.Unmarshal(self.Settings, &_settings)
		if err != nil {
			panic(err)
		}
		settings = _settings

		var _jobdata TrustpilotCollectorJobdata
		err = json.Unmarshal(self.Jobdata, &_jobdata)
		if err != nil {
			panic(err)
		}
		jobdata = _jobdata

	case CollectorTypePlayStore:
		var _settings PlayStoreCollectorSettings
		err := json.Unmarshal(self.Settings, &_settings)
		if err != nil {
			panic(err)
		}
		settings = _settings

		var _jobdata PlayStoreCollectorJobdata
		err = json.Unmarshal(self.Jobdata, &_jobdata)
		if err != nil {
			panic(err)
		}
		jobdata = _jobdata

	case CollectorTypeAppStore:
		var _settings AppStoreCollectorSettings
		err := json.Unmarshal(self.Settings, &_settings)
		settings = _settings
		if err != nil {
			panic(err)
		}

		var _jobdata AppStoreCollectorJobdata
		err = json.Unmarshal(self.Jobdata, &_jobdata)
		if err != nil {
			panic(err)
		}
		jobdata = _jobdata

	case CollectorTypeAmazon:
		var _settings AmazonCollectorSettings
		err := json.Unmarshal(self.Settings, &_settings)
		settings = _settings
		if err != nil {
			panic(err)
		}

		var _jobdata AmazonCollectorJobdata
		err = json.Unmarshal(self.Jobdata, &_jobdata)
		if err != nil {
			panic(err)
		}
		jobdata = _jobdata

	case CollectorTypeIAgora:
		var _settings IAgoraCollectorSettings
		err := json.Unmarshal(self.Settings, &_settings)
		if err != nil {
			panic(err)
		}
		settings = _settings

		var _jobdata IAgoraCollectorJobdata
		err = json.Unmarshal(self.Jobdata, &_jobdata)
		if err != nil {
			panic(err)
		}
		jobdata = _jobdata

	case CollectorTypeWebhook:
		var _settings WebhookCollectorSettings
		err := json.Unmarshal(self.Settings, &_settings)
		settings = _settings
		if err != nil {
			panic(err)
		}

		var _jobdata WebhookCollectorJobdata
		err = json.Unmarshal(self.Jobdata, &_jobdata)
		if err != nil {
			panic(err)
		}
		jobdata = _jobdata

	case CollectorTypeWidget:
		var _settings WidgetCollectorSettings
		err := json.Unmarshal(self.Settings, &_settings)
		settings = _settings
		if err != nil {
			panic(err)
		}

		var _jobdata WidgetCollectorJobdata
		err = json.Unmarshal(self.Jobdata, &_jobdata)
		if err != nil {
			panic(err)
		}
		jobdata = _jobdata

	default:
		panic(self.Type)
	}

	return &Collector{
		ID:        self.ID,
		ProductID: self.ProductID,
		Type:      self.Type,
		Settings:  settings,
		Jobdata:   jobdata,
		CreatedAt: self.CreatedAt,
		DeletedAt: self.DeletedAt,
	}
}
