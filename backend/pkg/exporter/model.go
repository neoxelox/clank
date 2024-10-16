package exporter

import (
	"encoding/json"
	"time"
)

const (
	EXPORTER_MODEL_TABLE = "\"exporter\""
)

type ExporterModel struct {
	ID        string     `db:"id"`
	ProductID string     `db:"product_id"`
	Type      string     `db:"type"`
	Settings  []byte     `db:"settings"`
	Jobdata   []byte     `db:"jobdata"`
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func NewExporterModel(exporter Exporter) *ExporterModel {
	settings, err := json.Marshal(exporter.Settings)
	if err != nil {
		panic(err)
	}

	jobdata, err := json.Marshal(exporter.Jobdata)
	if err != nil {
		panic(err)
	}

	return &ExporterModel{
		ID:        exporter.ID,
		ProductID: exporter.ProductID,
		Type:      exporter.Type,
		Settings:  settings,
		Jobdata:   jobdata,
		CreatedAt: exporter.CreatedAt,
		DeletedAt: exporter.DeletedAt,
	}
}

func (self *ExporterModel) ToEntity() *Exporter {
	var settings any
	var jobdata any
	switch self.Type {
	case ExporterTypeSlack:
		var _settings SlackExporterSettings
		err := json.Unmarshal(self.Settings, &_settings)
		if err != nil {
			panic(err)
		}
		settings = _settings

		var _jobdata SlackExporterJobdata
		err = json.Unmarshal(self.Jobdata, &_jobdata)
		if err != nil {
			panic(err)
		}
		jobdata = _jobdata

	case ExporterTypeJira:
		var _settings JiraExporterSettings
		err := json.Unmarshal(self.Settings, &_settings)
		if err != nil {
			panic(err)
		}
		settings = _settings

		var _jobdata JiraExporterJobdata
		err = json.Unmarshal(self.Jobdata, &_jobdata)
		if err != nil {
			panic(err)
		}
		jobdata = _jobdata

	default:
		panic(self.Type)
	}

	return &Exporter{
		ID:        self.ID,
		ProductID: self.ProductID,
		Type:      self.Type,
		Settings:  settings,
		Jobdata:   jobdata,
		CreatedAt: self.CreatedAt,
		DeletedAt: self.DeletedAt,
	}
}
